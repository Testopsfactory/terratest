package azure

import "testing"

func Test_getTargetAzureSubscription(t *testing.T) {
	type args struct {
		subID string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "subIDProvidedAsArg", args: args{subID: "test"}, want: "test", wantErr: false},
		{name: "subIDNotProvided", args: args{subID: ""}, want: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTargetAzureSubscription(tt.args.subID)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTargetAzureSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getTargetAzureSubscription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTargetAzureResourceGroupName(t *testing.T) {
	type args struct {
		rgName string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "rgNameProvidedAsArg", args: args{rgName: "test"}, want: "test", wantErr: false},
		{name: "rgNameNotProvided", args: args{rgName: ""}, want: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTargetAzureResourceGroupName(tt.args.rgName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTargetAzureResourceGroupName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getTargetAzureResourceGroupName() = %v, want %v", got, tt.want)
			}
		})
	}
}
