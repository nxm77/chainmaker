package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const (
	nodeId1      = "node1"
	nodeId2      = "node2"
	nodeAddress1 = "12345"
	nodeAddress2 = "22345"
)

func insertNodeTest() ([]*db.Node, error) {
	insertList := []*db.Node{
		{
			NodeId:  nodeId1,
			OrgId:   orgId1,
			Address: nodeAddress1,
		},
		{
			NodeId:  nodeId2,
			OrgId:   orgId2,
			Address: nodeAddress2,
		},
	}
	err := BatchInsertNode(ChainID, insertList)
	return insertList, err
}

func TestBatchInsertNode(t *testing.T) {
	insertList := []*db.Node{
		{
			NodeId:  nodeId1,
			OrgId:   orgId1,
			Address: nodeAddress1,
		},
		{
			NodeId:  nodeId2,
			OrgId:   orgId2,
			Address: nodeAddress2,
		},
	}

	type args struct {
		chainId  string
		nodeList []*db.Node
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:  ChainID,
				nodeList: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BatchInsertNode(tt.args.chainId, tt.args.nodeList); (err != nil) != tt.wantErr {
				t.Errorf("BatchInsertNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetNodeInById(t *testing.T) {
	insertList, err := insertNodeTest()
	if err != nil {
		return
	}
	type args struct {
		chainId string
		nodeIds []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*db.Node
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				nodeIds: []string{
					nodeId1,
				},
			},
			want: map[string]*db.Node{
				nodeId1: insertList[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNodeInById(tt.args.chainId, tt.args.nodeIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodeInById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Node{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetNodeInById() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetNodeList(t *testing.T) {
	insertList, err := insertNodeTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		orgId   string
		nodeId  string
		offset  int
		limit   int
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Node
		want1   int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				offset:  0,
				limit:   10,
			},
			want:    insertList,
			want1:   int64(len(insertList)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := GetNodeList(tt.args.chainId, tt.args.orgId, tt.args.nodeId, tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodeList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetNodeNum(t *testing.T) {
	insertList, err := insertNodeTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		role    string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want:    int64(len(insertList)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetNodeNum(tt.args.chainId, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodeNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetNodeNumByOrg(t *testing.T) {
	_, err := insertNodeTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		orgId   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId: db.UTchainID,
			},
			wantErr: false,
		},
		{
			name: "test: case 3",
			args: args{
				chainId: db.UTchainID,
				orgId:   "org1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetNodeNumByOrg(tt.args.chainId, tt.args.orgId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodeNumByOrg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetNodesRef(t *testing.T) {
	_, err := insertNodeTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want: []string{
				nodeId1,
				nodeId2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetNodesRef(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodesRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDeleteNodeById(t *testing.T) {
	err := DeleteNodeById(db.UTchainID, []string{"123"})
	if err != nil {
		t.Errorf("Test case 1 failed: %v", err)
	}
}
