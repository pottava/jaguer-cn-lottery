package googlecloud

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pottava/jaguer-cn-lottery/api/internal/lib"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func getService(ctx context.Context) (*sheets.Service, error) {
	if _, err := os.Stat(lib.Config.GcloudCreds); err == nil {
		bytes, err := ioutil.ReadFile(lib.Config.GcloudCreds)
		if err != nil {
			return nil, err
		}
		jwt, err := google.JWTConfigFromJSON(bytes, "https://www.googleapis.com/auth/spreadsheets")
		if err != nil {
			return nil, err
		}
		svc, err := sheets.NewService(ctx, option.WithHTTPClient(jwt.Client(ctx)))
		if err != nil {
			return nil, fmt.Errorf("unable to make a Sheets service: %v", err)
		}
		return svc, nil
	}
	svc, err := sheets.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to make a Sheets service: %v", err)
	}
	return svc, nil
}

func UpdateSheetCell(ctx context.Context, sheetID, tabID string, values []interface{}) error {
	svc, err := getService(ctx)
	if err != nil {
		return fmt.Errorf("unable to retrieve Sheets client: %v", err)
	}
	rangedValue := &sheets.ValueRange{Values: [][]interface{}{values}}
	if _, err = svc.Spreadsheets.Values.Append(sheetID, tabID, rangedValue).
		ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Do(); err != nil {
		return fmt.Errorf("unable to append a record to the sheet: %v", err)
	}
	return nil
}
