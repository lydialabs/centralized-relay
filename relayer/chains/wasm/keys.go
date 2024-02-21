package wasm

import (
	"context"
	"os"
	"path"
)

func (p *Provider) RestoreKeystore(ctx context.Context) error {
	filePath := path.Join(p.cfg.HomeDir, "keystore", p.NID(), p.cfg.Address)
	privFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	priv, err := p.kms.Decrypt(ctx, privFile)
	if err != nil {
		return err
	}
	passFile, err := os.ReadFile(filePath + ".pass")
	if err != nil {
		return err
	}
	pass, err := p.kms.Decrypt(ctx, passFile)
	if err != nil {
		return err
	}
	if err := p.client.ImportArmor(p.NID(), priv, string(pass)); err != nil {
		return err
	}
	return nil
}

func (p *Provider) NewKeystore(passphrase string) (string, error) {
	armor, addr, err := p.client.CreateAccount(p.NID(), passphrase)
	if err != nil {
		return "", err
	}
	encryptedArmor, err := p.kms.Encrypt(context.Background(), []byte(armor))
	if err != nil {
		return "", err
	}
	filePath := path.Join(p.cfg.HomeDir, "keystore", p.NID(), addr)
	if err = os.WriteFile(filePath, encryptedArmor, 0o644); err != nil {
		return "", err
	}
	encryptedPassphrase, err := p.kms.Encrypt(context.Background(), []byte(passphrase))
	if err != nil {
		return "", err
	}
	if err = os.WriteFile(filePath+".pass", encryptedPassphrase, 0o644); err != nil {
		return "", err
	}
	return addr, nil
}

// ImportKeystore imports a keystore from a file
func (p *Provider) ImportKeystore(ctx context.Context, keyPath, passphrase string) (string, error) {
	privFile, err := os.ReadFile(keyPath)
	if err != nil {
		return "", err
	}
	if err := p.client.ImportArmor(p.NID(), privFile, passphrase); err != nil {
		return "", err
	}
	armor, err := p.client.GetArmor(p.NID(), passphrase)
	if err != nil {
		return "", err
	}
	record, err := p.client.GetAccount(p.NID())
	if err != nil {
		return "", err
	}
	addr, err := record.GetAddress()
	if err != nil {
		return "", err
	}
	armorCipher, err := p.kms.Encrypt(ctx, []byte((armor)))
	if err != nil {
		return "", err
	}
	passphraseCipher, err := p.kms.Encrypt(ctx, []byte(passphrase))
	if err != nil {
		return "", err
	}
	filePath := path.Join(p.cfg.HomeDir, "keystore", p.NID(), addr.String())
	if err = os.WriteFile(filePath, armorCipher, 0o644); err != nil {
		return "", err
	}
	if err = os.WriteFile(filePath+".pass", passphraseCipher, 0o644); err != nil {
		return "", err
	}
	return addr.String(), nil
}
