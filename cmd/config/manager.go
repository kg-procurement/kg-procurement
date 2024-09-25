package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/tidwall/jsonc"
)

const (
	localConfigFileName = "config.jsonc"
)

type configManager struct {
	goValidator   *validator.Validate
	decoderConfig mapstructure.DecoderConfig
}

func (c *configManager) Start(ctx context.Context, dest any) error {
	err := c.readLocal(ctx, dest)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("config: local: %w", err)
		}
		fmt.Printf("Warning: local config file doesn't exist: %s\n", localConfigFileName)
	}

	return nil
}

func (c *configManager) readLocal(ctx context.Context, dest any) error {
	localConfRaw, err := os.ReadFile(localConfigFileName)
	if err != nil {
		return err
	}

	var localMap map[string]any
	localConfRaw = jsonc.ToJSON(localConfRaw)
	if err := json.Unmarshal(localConfRaw, &localMap); err != nil {
		return fmt.Errorf("failed to parse JSONC: %w", err)
	}

	if err := c.decodeAndValidate(ctx, localMap, dest); err != nil {
		return err
	}

	return nil
}

func (c *configManager) decodeAndValidate(ctx context.Context, input any, dest any) error {
	if err := c.decode(input, &dest); err != nil {
		return fmt.Errorf("failed to decode values into %T: %w", dest, err)
	}

	if err := c.goValidator.StructCtx(ctx, dest); err != nil {
		js, _ := json.Marshal(dest)
		return fmt.Errorf("failed to validate '%s' in %T: %w", string(js), dest, err)
	}
	return nil
}

func (c *configManager) decode(input any, dest any) (err error) {
	defer func() {
		if panicked := recover(); panicked != nil {
			err = fmt.Errorf("panicked when decoding config: %v", panicked)
		}
	}()
	dc := c.decoderConfig
	dc.Result = dest
	mdc, err := mapstructure.NewDecoder(&dc)
	if err != nil {
		return err
	}
	return mdc.Decode(input)
}

func NewConfigManager() *configManager {
	return &configManager{
		goValidator: validator.New(),
		decoderConfig: mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
			),
		},
	}
}
