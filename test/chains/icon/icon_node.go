package icon

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/icon-project/centralized-relay/test/interchaintest"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/blockdb"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/dockerutil"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/centralized-relay/test/interchaintest/relayer/centralized"
	"github.com/icon-project/centralized-relay/test/testsuite/testconfig"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon"
	"gopkg.in/yaml.v3"

	"github.com/icza/dyno"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	dockerclient "github.com/docker/docker/client"
	"github.com/icon-project/centralized-relay/test/chains"
	iconclient "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon"
	icontypes "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon/types"
	iconlog "github.com/icon-project/icon-bridge/common/log"
	"go.uber.org/zap"
)

const (
	rpcPort              = "9080/tcp"
	GOLOOP_IMAGE_ENV     = "GOLOOP_IMAGE"
	GOLOOP_IMAGE         = "iconloop/goloop-icon"
	GOLOOP_IMAGE_TAG_ENV = "GOLOOP_IMAGE_TAG"
	GOLOOP_IMAGE_TAG     = "latest"
)

var ContainerEnvs = [9]string{
	"GOCHAIN_CONFIG=/goloop/data/config.json",
	"GOCHAIN_GENESIS=/goloop/data/genesis.json",
	"GOCHAIN_DATA=/goloop/chain/iconee",
	"GOCHAIN_LOGFILE=/goloop/chain/iconee.log",
	"GOCHAIN_DB_TYPE=rocksdb",
	"GOCHAIN_CLEAN_DATA=true",
	"JAVAEE_BIN=/goloop/execman/bin/execman",
	"PYEE_VERIFY_PACKAGE=true",
	"ICON_CONFIG=/goloop/data/icon_config.json",
}

type IconRemoteNode struct {
	VolumeName   string
	Index        int
	Chain        chains.Chain
	NetworkID    string
	DockerClient *dockerclient.Client
	Client       icon.Client
	TestName     string
	Image        ibc.DockerImage
	log          *zap.Logger
	ContainerID  string
	// Ports set during StartContainer.
	HostRPCPort string
	Validator   bool
	lock        sync.Mutex
	Address     string
	testconfig  *testconfig.Chain
}

type IconNodes []*IconRemoteNode

// Name of the test node container
func (in *IconRemoteNode) Name() string {
	var nodeType string
	if in.Validator {
		nodeType = "val"
	} else {
		nodeType = "fn"
	}
	return fmt.Sprintf("%s-%s-%d-%s", in.Chain.Config().ChainID, nodeType, in.Index, dockerutil.SanitizeContainerName(in.TestName))
}

// Create Node Container with ports exposed and published for host to communicate with
func (in *IconRemoteNode) CreateNodeContainer(ctx context.Context, additionalGenesisWallets ...ibc.WalletAmount) error {
	imageRef := in.Image.Ref()
	testBasePath := os.Getenv(chains.BASE_PATH)
	binds := in.Bind()
	binds = append(binds, fmt.Sprintf("%s/test/chains/icon/data/governance:%s", testBasePath, "/goloop/data/gov"))
	containerConfig := &types.ContainerCreateConfig{
		Config: &container.Config{
			Image:    imageRef,
			Hostname: in.HostName(),
			Env:      ContainerEnvs[:],
			Labels:   map[string]string{dockerutil.CleanupLabel: in.TestName},
		},

		HostConfig: &container.HostConfig{
			Binds:           binds,
			PublishAllPorts: true,
			AutoRemove:      false,
			DNS:             []string{},
		},
		NetworkingConfig: &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				in.NetworkID: {},
			},
		},
	}
	cc, err := in.DockerClient.ContainerCreate(ctx, containerConfig.Config, containerConfig.HostConfig, containerConfig.NetworkingConfig, nil, in.Name())
	if err != nil {
		in.log.Error("Failed to create container", zap.Error(err))
		return err
	}

	err = in.modifyGenesisToAddGenesisAccount(ctx, cc.ID, additionalGenesisWallets...)
	if err != nil {
		in.log.Error("Failed to update genesis file in container", zap.Error(err))
		return err
	}

	in.CopyConfig(ctx, err, cc)

	in.ContainerID = cc.ID
	return nil
}

func (in *IconRemoteNode) CopyConfig(ctx context.Context, err error, cc container.CreateResponse) {
	fileName := fmt.Sprintf("%s/test/chains/icon/data/config.json", os.Getenv(chains.BASE_PATH))

	config, err := interchaintest.GetLocalFileContent(fileName)

	header := map[string]string{
		"name": "config.json",
	}
	err = in.CopyFileToContainer(context.WithValue(ctx, "file-header", header), config, cc.ID, "/goloop/data/")
}
func (in *IconRemoteNode) CopyFileToContainer(ctx context.Context, content []byte, containerID, target string) error {
	header := ctx.Value("file-header").(map[string]string)
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	err := tw.WriteHeader(&tar.Header{
		Name: header["name"],
		Mode: 0644,
		Size: int64(len(content)),
	})
	_, err = tw.Write(content)
	if err != nil {
		return err
	}
	err = tw.Close()
	if err != nil {
		return err
	}
	if err := in.DockerClient.CopyToContainer(context.Background(), containerID, target, &buf, types.CopyToContainerOptions{}); err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}

func (in *IconRemoteNode) modifyGenesisToAddGenesisAccount(ctx context.Context, containerID string, additionalGenesisWallets ...ibc.WalletAmount) error {
	g := make(map[string]interface{})
	fileName := fmt.Sprintf("%s/test/chains/icon/data/genesis.json", os.Getenv(chains.BASE_PATH))

	genbz, err := interchaintest.GetLocalFileContent(fileName)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(genbz, &g); err != nil {
		return fmt.Errorf("failed to unmarshal genesis file: %w", err)
	}

	for index, wallet := range additionalGenesisWallets {
		genesisAccount := map[string]string{
			"address": wallet.Address,
			"balance": "0xd3c21bcecceda1000000", // 1_000_000*10**18
			"name":    fmt.Sprintf("ibc-%d", index),
		}
		if err := dyno.Append(g, genesisAccount, "accounts"); err != nil {
			return fmt.Errorf("failed to set add genesis accounts in genesis json: %w", err)
		}
		in.log.Info("Genesis file update with faucet wallet",
			zap.String("wallet", wallet.Address),
			zap.String("amount", "0xd3c21bcecceda1000000"),
		)
	}
	result, _ := json.Marshal(g)

	header := map[string]string{
		"name": "genesis.json",
	}
	err = in.CopyFileToContainer(context.WithValue(ctx, "file-header", header), result, containerID, "/goloop/data/")

	return err
}

func (in *IconRemoteNode) HostName() string {
	return dockerutil.CondenseHostName(in.Name())
}

func (in *IconRemoteNode) Bind() []string {
	return []string{fmt.Sprintf("%s:%s", in.VolumeName, in.HomeDir())}
}

func (in *IconRemoteNode) HomeDir() string {
	return path.Join("/var/icon-chain", in.Chain.Config().Name)
}

func (in *IconRemoteNode) StartContainer(ctx context.Context) error {
	if err := dockerutil.StartContainer(ctx, in.DockerClient, in.ContainerID); err != nil {
		return err
	}

	c, err := in.DockerClient.ContainerInspect(ctx, in.ContainerID)
	if err != nil {
		return err
	}
	in.HostRPCPort = dockerutil.GetHostPort(c, rpcPort)
	in.logger().Info("Icon chain node started", zap.String("container", in.Name()), zap.String("rpc_port", in.HostRPCPort))

	uri := "http://" + in.HostRPCPort + "/api/v3/"
	var l iconlog.Logger
	in.Client = *iconclient.NewClient(uri, l)
	return nil
}

func (in *IconRemoteNode) logger() *zap.Logger {
	return in.log.With(
		zap.String("chain_id", in.Chain.Config().ChainID),
		zap.String("test", in.TestName),
	)
}

func (in *IconRemoteNode) Exec(ctx context.Context, cmd []string, env []string) ([]byte, []byte, error) {
	job := dockerutil.NewImage(in.logger(), in.DockerClient, in.NetworkID, in.TestName, chains.GetEnvOrDefault(GOLOOP_IMAGE_ENV, GOLOOP_IMAGE), chains.GetEnvOrDefault(GOLOOP_IMAGE_TAG_ENV, GOLOOP_IMAGE_TAG))
	opts := dockerutil.ContainerOptions{
		Binds: []string{
			in.testconfig.ContractsPath + ":/contracts",
			in.testconfig.ConfigPath + ":/goloop/data",
		},
		Env: ContainerEnvs[:],
	}
	// opts := dockerutil.ContainerOptions{
	// 	Env:   env,
	// 	Binds: in.Bind(),
	// }
	res := job.Run(ctx, cmd, opts)
	return res.Stdout, res.Stderr, res.Err
}

func (in *IconRemoteNode) BinCommand(command ...string) []string {
	command = append([]string{in.Chain.Config().Bin}, command...)
	return command
}

func (in *IconRemoteNode) ExecBin(ctx context.Context, command ...string) ([]byte, []byte, error) {
	return in.Exec(ctx, in.BinCommand(command...), nil)
}

func (in *IconRemoteNode) GetBlockByHeight(ctx context.Context, height int64) (string, error) {
	in.lock.Lock()
	defer in.lock.Unlock()
	uri := "http://" + in.HostRPCPort + "/api/v3"
	block, _, err := in.ExecBin(ctx,
		"rpc", "blockbyheight", fmt.Sprint(height),
		"--uri", uri,
	)
	return string(block), err
}

func (in *IconRemoteNode) FindTxs(ctx context.Context, height uint64) ([]blockdb.Tx, error) {
	var flag = true
	if flag {
		time.Sleep(3 * time.Second)
		flag = false
	}

	time.Sleep(2 * time.Second)
	blockHeight := icontypes.BlockHeightParam{Height: icontypes.NewHexInt(int64(height))}
	res, err := in.Client.GetBlockByHeight(&blockHeight)
	if err != nil {
		return make([]blockdb.Tx, 0, 0), nil
	}
	txs := make([]blockdb.Tx, 0, len(res.NormalTransactions)+2)
	var newTx blockdb.Tx
	for _, tx := range res.NormalTransactions {
		newTx.Data = []byte(fmt.Sprintf(`{"data":"%s"}`, tx.Data))
	}

	// ToDo Add events from block if any to newTx.Events.
	// Event is an alternative representation of tendermint/abci/types.Event
	return txs, nil
}

func (in *IconRemoteNode) Height(ctx context.Context) (uint64, error) {
	res, err := in.Client.GetLastBlock()
	return uint64(res.Height), err
}

func (in *IconRemoteNode) GetBalance(ctx context.Context, address string) (int64, error) {
	addr := icontypes.AddressParam{Address: icontypes.Address(address)}
	bal, err := in.Client.GetBalance(&addr)
	return bal.Int64(), err
}

func (in *IconRemoteNode) DeployContract(ctx context.Context, scorePath, keystorePath, initMessage string) (string, error) {
	// Write Contract file to Docker volume
	_, score := filepath.Split(scorePath)
	if err := in.CopyFile(ctx, scorePath, score); err != nil {
		return "", fmt.Errorf("error copying 6123f953784d27e0729bc7a640d6ad8f04ed6710.keystore to Docker volume: %w", err)
	}

	// Deploy the contract
	hash, err := in.ExecTx(ctx, initMessage, path.Join(in.HomeDir(), score), keystorePath)
	if err != nil {
		return "", err
	}

	//wait for few blocks
	time.Sleep(3 * time.Second)

	// Get Score Address
	trResult, err := in.TransactionResult(ctx, hash)

	if err != nil {
		return "", err
	}

	return string(trResult.SCOREAddress), nil

}

// Get Transaction result when hash is provided after executing a transaction
func (in *IconRemoteNode) TransactionResult(ctx context.Context, hash string) (*icontypes.TransactionResult, error) {
	uri := fmt.Sprintf("http://%s:9080/api/v3", in.Name()) //"http://" + in.HostRPCPort + "/api/v3"
	out, _, err := in.ExecBin(ctx, "rpc", "txresult", hash, "--uri", uri)
	if err != nil {
		return nil, err
	}
	var result = new(icontypes.TransactionResult)
	return result, json.Unmarshal(out, result)
}

// ExecTx executes a transaction, waits for 2 blocks if successful, then returns the tx hash.
func (in *IconRemoteNode) ExecTx(ctx context.Context, initMessage string, filePath string, keystorePath string, command ...string) (string, error) {
	var output string
	in.lock.Lock()
	defer in.lock.Unlock()
	stdout, _, err := in.Exec(ctx, in.TxCommand(ctx, initMessage, filePath, keystorePath, command...), nil)
	if err != nil {
		return "", err
	}
	return output, json.Unmarshal(stdout, &output)
}

// TxCommand is a helper to retrieve a full command for broadcasting a tx
// with the chain node binary.
func (in *IconRemoteNode) TxCommand(ctx context.Context, initMessage, filePath, keystorePath string, command ...string) []string {
	// get password from pathname as pathname will have the password prefixed. ex - Alice.Json
	_, key := filepath.Split(keystorePath)
	fileName := strings.Split(key, ".")
	password := fileName[0]

	command = append([]string{"rpc", "sendtx", "deploy", filePath}, command...)
	command = append(command,
		"--key_store", keystorePath,
		"--key_password", password,
		"--step_limit", "5000000000",
		"--content_type", "application/java",
	)
	if initMessage != "" && initMessage != "{}" {
		if strings.HasPrefix(initMessage, "{") {
			command = append(command, "--params", initMessage)
		} else {
			command = append(command, "--param", initMessage)
		}
	}

	return in.NodeCommand(command...)
}

// NodeCommand is a helper to retrieve a full command for a chain node binary.
// when interactions with the RPC endpoint are necessary.
// For example, if chain node binary is `gaiad`, and desired command is `gaiad keys show key1`,
// pass ("keys", "show", "key1") for command to return the full command.
// Will include additional flags for node URL, home directory, and chain ID.
func (in *IconRemoteNode) NodeCommand(command ...string) []string {
	command = in.BinCommand(command...)
	return append(command,
		"--uri", fmt.Sprintf("http://%s:9080/api/v3", in.Name()), //fmt.Sprintf("http://%s/api/v3", in.HostRPCPort),
		"--nid", "0x3",
	)
}

// CopyFile adds a file from the host filesystem to the docker filesystem
// relPath describes the location of the file in the docker volume relative to
// the home directory
func (tn *IconRemoteNode) CopyFile(ctx context.Context, srcPath, dstPath string) error {
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	return tn.WriteFile(ctx, content, dstPath)
}

// WriteFile accepts file contents in a byte slice and writes the contents to
// the docker filesystem. relPath describes the location of the file in the
// docker volume relative to the home directory
func (tn *IconRemoteNode) WriteFile(ctx context.Context, content []byte, relPath string) error {
	fw := dockerutil.NewFileWriter(tn.logger(), tn.DockerClient, tn.TestName)
	return fw.WriteFile(ctx, tn.VolumeName, relPath, content)
}

func (in *IconRemoteNode) QueryContract(ctx context.Context, scoreAddress, methodName, params string) ([]byte, error) {
	uri := fmt.Sprintf("http://%s:9080/api/v3", in.Name())
	var args = []string{"rpc", "call", "--to", scoreAddress, "--method", methodName, "--uri", uri}
	if params != "" {
		var paramName = "--param"
		if strings.HasPrefix(params, "{") && strings.HasSuffix(params, "}") {
			paramName = "--raw"
		}
		args = append(args, paramName, params)
	}
	out, _, err := in.ExecBin(ctx, args...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (in *IconRemoteNode) RestoreKeystore(ctx context.Context, ks []byte, keyName string) error {
	return in.WriteFile(ctx, ks, keyName+".json")
}

func (in *IconRemoteNode) ExecuteContract(ctx context.Context, scoreAddress, methodName, keyStorePath, params string) (string, error) {
	return in.ExecCallTx(ctx, scoreAddress, methodName, keyStorePath, params)
}

func (in *IconRemoteNode) ExecCallTx(ctx context.Context, scoreAddress, methodName, keystorePath, params string) (string, error) {
	var output string
	in.lock.Lock()
	defer in.lock.Unlock()
	stdout, _, err := in.Exec(ctx, in.ExecCallTxCommand(ctx, scoreAddress, methodName, keystorePath, params), nil)
	if err != nil {
		return "", err
	}
	return output, json.Unmarshal(stdout, &output)
}

func (in *IconRemoteNode) ExecCallTxCommand(ctx context.Context, scoreAddress, methodName, keystorePath, params string) []string {
	// get password from pathname as pathname will have the password prefixed. ex - Alice.Json
	_, key := filepath.Split(keystorePath)
	fileName := strings.Split(key, ".")
	password := fileName[0]
	command := []string{"rpc", "sendtx", "call"}

	command = append(command,
		"--to", scoreAddress,
		"--method", methodName,
		"--key_store", keystorePath,
		"--key_password", password,
		"--step_limit", "5000000000",
	)

	if params != "" && params != "{}" {
		if strings.HasPrefix(params, "{") {
			command = append(command, "--params", params)
		} else {
			command = append(command, "--param", params)
		}
	}

	if methodName == "registerPRep" {
		command = append(command, "--value", "2000000000000000000000")
	}

	return in.NodeCommand(command...)
}

func (in *IconRemoteNode) GetDebugTrace(ctx context.Context, hash icontypes.HexBytes) (*DebugTrace, error) {
	uri := fmt.Sprintf("http://%s:9080/api/v3d", in.Name())
	out, _, err := in.ExecBin(ctx, "debug", "trace", string(hash), "--uri", uri)
	if err != nil {
		return nil, err
	}
	var result = new(DebugTrace)
	return result, json.Unmarshal(out, result)

}

func (in *IconRemoteNode) GetChainConfig(ctx context.Context, rlyHome string, keyName string) ([]byte, error) {

	config := &centralized.ICONRelayerChainConfig{
		Type: "icon",
		Value: centralized.ICONRelayerChainConfigValue{
			NID:         in.Chain.Config().ChainID,
			RPCURL:      in.Chain.GetRPCAddress(),
			StartHeight: 0,
			NetworkID:   0x3,
		},
	}
	return yaml.Marshal(config)
}
