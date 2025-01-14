package configrpc

import (
	"context"
	"reflect"
	"testing"

	"github.com/mcku/retrygger/v2/modules/grpc/reconpb/jobmgmt"
)

func TestReadConfigRpc(t *testing.T) {
	type args struct {
		ctx          context.Context
		providerAddr string
		service      string
		job          string
	}
	tests := []struct {
		name    string
		args    args
		want    *jobmgmt.JobConfig
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				ctx:          context.TODO(),
				providerAddr: "127.0.0.1:18090",
				service:      "retrygger-test",
				job:          "mikrosaray-pos-cron",
			},
			want: &jobmgmt.JobConfig{
				RecordId:       "retrygger-mikrosaray-pos-cron",
				Service:        "retrygger-test",
				Job:            "mikrosaray-pos-cron",
				CronSchedule:   "0 */3 * * * *",
				ConfigDbParams: "",
				Enabled:        jobmgmt.YesNoStatus_YES_NO_STATUS_NO,
				SendNotif:      jobmgmt.YesNoStatus_YES_NO_STATUS_YES,
				CreatedAt:      1732267589118000000,
				UpdatedAt:      1732267589118000000,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadConfigRpc(tt.args.ctx, tt.args.providerAddr, tt.args.service, tt.args.job)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfigRpc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfigRpc() = %v, want %v", got, tt.want)
			}
		})
	}
}
