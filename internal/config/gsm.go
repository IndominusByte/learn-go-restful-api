package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/encryption"
)

type gsmData struct {
	EncryptionKey  string `json:"encryption_key"`
	PgTalkUser     string `json:"pg_talk_user"`
	PgTalkPassword string `json:"pg_talk_password"`
}

func (cfg *Config) loadFromGsm() error {
	filename := fmt.Sprintf("./conf/data.%s.json", os.Getenv("BACKEND_STAGE"))
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var data gsmData
	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		return err
	}

	// decode data
	cdn := &encryption.Credentials{Key: []byte(data.EncryptionKey)}

	pgtalkuser, pgtalkusererr := cdn.Decrypt(data.PgTalkUser)
	if pgtalkusererr != nil {
		return err
	}

	pgtalkpassword, pgtalkpassworderr := cdn.Decrypt(data.PgTalkPassword)
	if pgtalkpassworderr != nil {
		return pgtalkpassworderr
	}

	cfg.Database.MasterDsn = fmt.Sprintf(cfg.Database.MasterDsnNoCred, pgtalkuser, pgtalkpassword)
	cfg.Database.FollowerDsn = fmt.Sprintf(cfg.Database.FollowerDsnNoCred, pgtalkuser, pgtalkpassword)

	return nil
}
