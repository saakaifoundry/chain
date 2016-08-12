package core

import (
	"golang.org/x/net/context"

	"chain/core/txbuilder"
	"chain/crypto/ed25519/hd25519"
	"chain/encoding/json"
	"chain/errors"
	"chain/net/http/httpjson"
)

func (a *api) mockhsmCreateKey(ctx context.Context) (result struct {
	XPub json.HexBytes `json:"xpub"`
}, err error) {
	xpub, err := a.hsm.CreateKey(ctx)
	if err != nil {
		return result, err
	}
	result.XPub = xpub.Bytes()
	return result, nil
}

func (a *api) mockhsmListKeys(ctx context.Context, query struct{ Cursor string }) (page, error) {
	limit := defGenericPageSize

	xpubs, cursor, err := a.hsm.ListKeys(ctx, query.Cursor, limit)
	if err != nil {
		return page{}, err
	}

	var items []interface{}
	for _, xpub := range xpubs {
		item := struct {
			XPub *hd25519.XPub `json:"xpub"`
		}{
			xpub,
		}
		items = append(items, item)
	}

	return page{
		Items:    httpjson.Array(items),
		LastPage: len(xpubs) < limit,
		Query:    requestQuery{Cursor: cursor},
	}, nil
}

func (a *api) mockhsmDelKey(ctx context.Context, xpubBytes json.HexBytes) error {
	xpub, err := hd25519.XPubFromBytes(xpubBytes)
	if err != nil {
		return err
	}
	return a.hsm.DelKey(ctx, xpub)
}

func (a *api) mockhsmSignTemplates(ctx context.Context, tpls []*txbuilder.Template) []interface{} {
	resp := make([]interface{}, 0, len(tpls))
	for _, tpl := range tpls {
		err := txbuilder.Sign(ctx, tpl, a.mockhsmSignTemplate)
		if err != nil {
			info, _ := errInfo(err)
			resp = append(resp, info)
		} else {
			resp = append(resp, tpl)
		}
	}
	return resp
}

func (a *api) mockhsmSignTemplate(ctx context.Context, sigComponent *txbuilder.SigScriptComponent, sig *txbuilder.Signature) ([]byte, error) {
	xpub, err := hd25519.XPubFromString(sig.XPub)
	if err != nil {
		return nil, errors.Wrap(err, "parsing xpub")
	}
	return a.hsm.Sign(ctx, xpub, sig.DerivationPath, sigComponent.SignatureData[:])
}