package task

import (
	"context"
	"fmt"
	"time"

	paycli "github.com/hardstylez72/cry-pay/proto/gen/go/v1"
	"github.com/hardstylez72/cry/internal/orbiter"
	"github.com/hardstylez72/cry/internal/pay"
	"github.com/hardstylez72/cry/internal/pb/gen/proto/go/v1"
	"github.com/hardstylez72/cry/internal/process/halp"
	"github.com/hardstylez72/cry/internal/server/config"
	"github.com/hardstylez72/cry/internal/server/repository"
	"github.com/hardstylez72/cry/internal/settings"
	"github.com/hardstylez72/cry/internal/snapshot"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Input struct {
	ProfileId string
	ProcessId string
	UserId    string
	Task      *v1.ProcessTask

	L                    *zap.SugaredLogger
	ProcessRepository    repository.ProcessRepository
	ProfileRepository    repository.ProfileRepository
	WithdrawerRepository repository.WithdrawerRepository
	SettingsService      *settings.Service

	Snapshot snapshot.Voter

	Halper     *halp.Halp
	PayService *pay.Service
	User       *repository.User
	Orbiter    *orbiter.Service
}

type TaskRes struct {
	Task *v1.ProcessTask
}

type EstimateSyncSwapCostReq struct {
	Task      *v1.ProcessTask
	UserId    string
	ProcessId string
	ProfileId string
}

type Tasker interface {
	Run(ctx context.Context, arg *Input) (*v1.ProcessTask, error)
	Stop() error
	Type() v1.TaskType
}

const (
	taskTimeout         = time.Minute * 2
	taskStarkNetTimeout = time.Minute * 10
)

var PayableTasks = []v1.TaskType{
	v1.TaskType_SyncSwap,
	v1.TaskType_StargateBridge,
	v1.TaskType_ZkSyncOfficialBridgeToEthereum,
	v1.TaskType_TestNetBridgeSwap,
	v1.TaskType_SnapshotVote,
	v1.TaskType_OrbiterBridge,
	v1.TaskType_ZkSyncOfficialBridgeFromEthereum,
	v1.TaskType_WETH,
	v1.TaskType_MuteioSwap,
	v1.TaskType_SyncSwapLP,
	v1.TaskType_MaverickSwap,
	v1.TaskType_SpaceFISwap,
	v1.TaskType_VelocoreSwap,
	v1.TaskType_IzumiSwap,
	v1.TaskType_VeSyncSwap,
	v1.TaskType_EzkaliburSwap,
	v1.TaskType_ZkSwap,
	v1.TaskType_TraderJoeSwap,
	v1.TaskType_MerklyMintAndBridgeNFT,
	v1.TaskType_DeployStarkNetAccount,
	v1.TaskType_Swap10k,
	v1.TaskType_PancakeSwap,
	v1.TaskType_SithSwap,
	v1.TaskType_JediSwap,
	v1.TaskType_MySwap,
	v1.TaskType_ProtossSwap,
	v1.TaskType_StarkNetBridge,
}

var NonPayableTasks = []v1.TaskType{
	v1.TaskType_Delay,
	v1.TaskType_OkexDeposit,
	v1.TaskType_WithdrawExchange,
}

var executors = map[v1.TaskType]Tasker{
	v1.TaskType_StargateBridge:                   &Wrap{Tasker: &StargateTask{}},
	v1.TaskType_Mock:                             &Wrap{Tasker: &mockTask{}},
	v1.TaskType_Delay:                            &Wrap{Tasker: &taskDelay{}},
	v1.TaskType_WithdrawExchange:                 &Wrap{Tasker: &WithdrawExchange{}},
	v1.TaskType_OkexDeposit:                      &Wrap{Tasker: &OkexDepositTask{}},
	v1.TaskType_TestNetBridgeSwap:                &Wrap{Tasker: &TestNetBridgeSwapTask{}},
	v1.TaskType_SnapshotVote:                     &Wrap{Tasker: &SnapshotVoteTask{}},
	v1.TaskType_OkexBinance:                      &Wrap{Tasker: &mockTask{}},
	v1.TaskType_SyncSwap:                         &Wrap{Tasker: NewSyncSwapTask()},
	v1.TaskType_ZkSyncOfficialBridgeToEthereum:   &Wrap{Tasker: &ZksyncOfficialBridgeToEthereumTask{}},
	v1.TaskType_OrbiterBridge:                    &Wrap{Tasker: &OrbiterBridgeTask{}},
	v1.TaskType_ZkSyncOfficialBridgeFromEthereum: &Wrap{Tasker: &ZksyncOfficialBridgeFromEthereumTask{}},
	v1.TaskType_WETH:                             &Wrap{Tasker: &WethTask{}},
	v1.TaskType_MuteioSwap:                       &Wrap{Tasker: NewMuteioSwapTask()},
	v1.TaskType_SyncSwapLP:                       &Wrap{Tasker: &SyncSwapLPTask{}},
	v1.TaskType_MaverickSwap:                     &Wrap{Tasker: NewMaverickSwapTask()},
	v1.TaskType_SpaceFISwap:                      &Wrap{Tasker: NewSpaceFiSwapTask()},
	v1.TaskType_VelocoreSwap:                     &Wrap{Tasker: NewVelocoreSwapTask()},
	v1.TaskType_IzumiSwap:                        &Wrap{Tasker: NewIzumiSwapTask()},
	v1.TaskType_VeSyncSwap:                       &Wrap{Tasker: NewVeSyncSwapTask()},
	v1.TaskType_EzkaliburSwap:                    &Wrap{Tasker: NewEzkaliburSwapTask()},
	v1.TaskType_ZkSwap:                           &Wrap{Tasker: NewZkSwapTask()},
	v1.TaskType_TraderJoeSwap:                    &Wrap{Tasker: NewTraderJoeSwapTask()},
	v1.TaskType_MerklyMintAndBridgeNFT:           &Wrap{Tasker: &MerklyMintAndBridgeNFTTask{}},
	v1.TaskType_DeployStarkNetAccount:            &Wrap{Tasker: &DeployStarkNetAccountTask{}},
	v1.TaskType_Swap10k:                          &Wrap{Tasker: NewSwap10kSwapTask()},
	v1.TaskType_PancakeSwap:                      &Wrap{Tasker: NewPancakeSwapTask()},
	v1.TaskType_SithSwap:                         &Wrap{Tasker: NewSithSwapTask()},
	v1.TaskType_JediSwap:                         &Wrap{Tasker: NewJediSwapTask()},
	v1.TaskType_MySwap:                           &Wrap{Tasker: NewMySwapSwapTask()},
	v1.TaskType_ProtossSwap:                      &Wrap{Tasker: NewProtossSwapTask()},
	v1.TaskType_StarkNetBridge:                   &Wrap{Tasker: NewStarkNetBridgeTask()},
}

func GetTaskDesc(m *v1.Task) ([]byte, error) {
	switch m.TaskType {
	case v1.TaskType_Delay:
		t, ok := m.Task.(*v1.Task_DelayTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_DelayTask)")
		}
		return Marshal(t.DelayTask)
	case v1.TaskType_SnapshotVote:
		t, ok := m.Task.(*v1.Task_SnapshotVoteTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_SnapshotVoteTask)")
		}
		return Marshal(t.SnapshotVoteTask)
	case v1.TaskType_StargateBridge:
		t, ok := m.Task.(*v1.Task_StargateBridgeTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_StargateBridgeTask)")
		}
		return Marshal(t.StargateBridgeTask)
	case v1.TaskType_Mock:
		t, ok := m.Task.(*v1.Task_MockTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_MockTask)")
		}
		return Marshal(t.MockTask)
	case v1.TaskType_TestNetBridgeSwap:
		t, ok := m.Task.(*v1.Task_TestNetBridgeSwapTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_TestNetBridgeSwapTask)")
		}
		return Marshal(t.TestNetBridgeSwapTask)
	case v1.TaskType_WithdrawExchange:
		t, ok := m.Task.(*v1.Task_WithdrawExchangeTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_WithdrawExchangeTask)")
		}
		return Marshal(t.WithdrawExchangeTask)
	case v1.TaskType_OkexDeposit:
		t, ok := m.Task.(*v1.Task_OkexDepositTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_OkexDepositTask)")
		}
		return Marshal(t.OkexDepositTask)
	case v1.TaskType_OkexBinance:
		t, ok := m.Task.(*v1.Task_OkexBinanaceTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_OkexDepositTask)")
		}
		return Marshal(t.OkexBinanaceTask)
	case v1.TaskType_SyncSwap:
		t, ok := m.Task.(*v1.Task_SyncSwapTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_SyncSwapTask)")
		}
		return Marshal(t.SyncSwapTask)
	case v1.TaskType_ZkSyncOfficialBridgeToEthereum:
		t, ok := m.Task.(*v1.Task_ZkSyncOfficialBridgeToEthereumTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_ZkSyncOfficialBridgeToEthereumTask)")
		}
		return Marshal(t.ZkSyncOfficialBridgeToEthereumTask)
	case v1.TaskType_OrbiterBridge:
		t, ok := m.Task.(*v1.Task_OrbiterBridgeTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_OrbiterBridgeTask)")
		}
		return Marshal(t.OrbiterBridgeTask)
	case v1.TaskType_ZkSyncOfficialBridgeFromEthereum:
		t, ok := m.Task.(*v1.Task_ZkSyncOfficialBridgeFromEthereumTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_ZkSyncOfficialBridgeFromEthereumTask)")
		}
		return Marshal(t.ZkSyncOfficialBridgeFromEthereumTask)

	case v1.TaskType_WETH:
		t, ok := m.Task.(*v1.Task_WETHTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_WETHTask)")
		}
		return Marshal(t.WETHTask)
	case v1.TaskType_MuteioSwap:
		t, ok := m.Task.(*v1.Task_MuteioSwapTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_MuteioSwapTask)")
		}
		return Marshal(t.MuteioSwapTask)
	case v1.TaskType_SyncSwapLP:
		t, ok := m.Task.(*v1.Task_SyncSwapLPTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_SyncSwapLPTask)")
		}
		return Marshal(t.SyncSwapLPTask)
	case v1.TaskType_MaverickSwap:
		t, ok := m.Task.(*v1.Task_MaverickSwapTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_MaverickSwapTask)")
		}
		return Marshal(t.MaverickSwapTask)
	case v1.TaskType_SpaceFISwap:
		t, ok := m.Task.(*v1.Task_SpaceFiSwapTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_SpaceFiSwapTask)")
		}
		return Marshal(t.SpaceFiSwapTask)
	case v1.TaskType_VelocoreSwap:
		t, ok := m.Task.(*v1.Task_VelocoreSwapTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_VelocoreSwapTask)")
		}
		return Marshal(t.VelocoreSwapTask)
	case v1.TaskType_IzumiSwap:
		t, ok := m.Task.(*v1.Task_IzumiSwapTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_IzumiSwapTask)")
		}
		return Marshal(t.IzumiSwapTask)
	case v1.TaskType_VeSyncSwap:
		t, ok := m.Task.(*v1.Task_VeSyncSwapTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_VeSyncSwapTask)")
		}
		return Marshal(t.VeSyncSwapTask)
	case v1.TaskType_EzkaliburSwap:
		t, ok := m.Task.(*v1.Task_EzkaliburSwapTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_EzkaliburSwapTask)")
		}
		return Marshal(t.EzkaliburSwapTask)
	case v1.TaskType_ZkSwap:
		t, ok := m.Task.(*v1.Task_ZkSwapTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_ZkSwapTask)")
		}
		return Marshal(t.ZkSwapTask)
	case v1.TaskType_TraderJoeSwap:
		t, ok := m.Task.(*v1.Task_TraderJoeSwapTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_TraderJoeSwapTask)")
		}
		return Marshal(t.TraderJoeSwapTask)
	case v1.TaskType_MerklyMintAndBridgeNFT:
		t, ok := m.Task.(*v1.Task_MerklyMintAndBridgeNFTTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_MerklyMintAndBridgeNFTTask)")
		}
		return Marshal(t.MerklyMintAndBridgeNFTTask)
	case v1.TaskType_DeployStarkNetAccount:
		t, ok := m.Task.(*v1.Task_DeployStarkNetAccountTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_DeployStarkNetAccountTask)")
		}
		return Marshal(t.DeployStarkNetAccountTask)
	case v1.TaskType_Swap10k:
		t, ok := m.Task.(*v1.Task_Swap10K)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_Swap10K)")
		}
		return Marshal(t.Swap10K)
	case v1.TaskType_PancakeSwap:
		t, ok := m.Task.(*v1.Task_PancakeSwapTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_PancakeSwapTask)")
		}
		return Marshal(t.PancakeSwapTask)

	case v1.TaskType_SithSwap:
		t, ok := m.Task.(*v1.Task_SithSwapTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_SithSwapTask)")
		}
		return Marshal(t.SithSwapTask)
	case v1.TaskType_JediSwap:
		t, ok := m.Task.(*v1.Task_JediSwapTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_JediSwapTask)")
		}
		return Marshal(t.JediSwapTask)
	case v1.TaskType_MySwap:
		t, ok := m.Task.(*v1.Task_MySwapTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_MySwapTask)")
		}
		return Marshal(t.MySwapTask)
	case v1.TaskType_ProtossSwap:
		t, ok := m.Task.(*v1.Task_ProtosSwapTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_ProtosSwapTask)")
		}
		return Marshal(t.ProtosSwapTask)
	case v1.TaskType_StarkNetBridge:
		t, ok := m.Task.(*v1.Task_StarkNetBridgeTask)
		if !ok {
			return nil, errors.New("m.Task.(*v1.Task_StarkNetBridgeTask)")
		}
		return Marshal(t.StarkNetBridgeTask)
	default:
		return nil, errors.New("invalid task type: " + m.TaskType.String())
	}
}

func GetTask(t v1.TaskType) (Tasker, error) {
	tasker, exist := executors[t]
	if !exist {
		return nil, errors.New("unknown task: " + t.String())
	}
	return tasker, nil
}

type Wrap struct {
	Tasker Tasker
}

func (w *Wrap) Stop() error {
	return w.Tasker.Stop()
}
func (w *Wrap) Type() v1.TaskType {
	return w.Tasker.Type()
}

type TaskUpdater interface {
	UpdateTask(ctx context.Context, req *v1.ProcessTask) error
}

func (d *Input) UpdateTask(ctx context.Context, req *v1.ProcessTask) error {
	return UpdateTask(ctx, req, d.ProcessRepository, d.PayService, d.User, d.UserId)
}

func UpdateTask(ctx context.Context, after *v1.ProcessTask, d repository.ProcessRepository, payService *pay.Service, user *repository.User, userId string) error {

	t, err := d.GetProcessTask(ctx, after.Id)
	if err != nil {
		return err
	}

	before, err := t.ToPB()
	if err != nil {
		return err
	}

	if after.Status != before.Status {
		if after.Error != nil {
			_ = d.RecordStatusChanged(ctx, before.Id, before.Status, after.Status, *after.Error)
		} else {
			_ = d.RecordStatusChanged(ctx, before.Id, before.Status, after.Status)
		}
	}

	marshal, err := GetTaskDesc(after.Task)
	if err != nil {
		return errors.Wrap(err, "GetTaskDesc")
	}
	after.Task.Description = string(marshal)

	if err := d.UpdateProcessTask(ctx, after, before.Id, t.ProcessId, t.ProfileId); err != nil {
		return err
	}

	if NeedPay(before, after) {
		_, err = payService.FundsServiceClient.TaskCompleted(ctx, &paycli.TaskCompletedReq{
			ProcessId: t.ProcessId,
			ProfileId: t.ProfileId,
			TaskId:    before.Id,
			TaskType:  before.Task.TaskType.String(),
			UserId:    userId,
		})
		if err != nil {
			return err
		}

	}

	return nil
}

func (w *Wrap) Run(ctx context.Context, a *Input) (task *v1.ProcessTask, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("task panic: %+v", r))
			task = nil
		}
	}()

	if IsPayableTask(a.Task.Task.TaskType) {
		funds, err := a.PayService.FundsServiceClient.GetFunds(ctx, &paycli.GetFundsReq{
			Login: a.User.Email,
		})
		if err != nil {
			return nil, err
		}

		if funds.FundsLeft <= 0 {
			e := ErrUserHasNoBalance.Error()
			a.Task.Error = &e
			a.Task.Status = v1.ProcessStatus_StatusError

			if err := a.UpdateTask(ctx, a.Task); err != nil {
				return nil, err
			}

			return nil, ErrUserHasNoBalance
		}
	}

	taskId := a.Task.Id

	l := a.L.With("taskId", taskId).
		With("task", a.Task.Task.TaskType.String()).
		With("user_id", a.UserId)

	tr := otel.Tracer("")
	pctx, span := tr.Start(ctx, "RunTask")
	span.SetAttributes(
		attribute.String("tId", taskId),
		attribute.String("pId", a.ProcessId),
		attribute.String("ppId", a.ProfileId),
		attribute.String("tType", a.Task.Task.TaskType.String()),
		attribute.String("userId", a.UserId),
	)
	defer span.End()

	l.Debug("task running")
	task, err = w.Tasker.Run(pctx, a)
	if err != nil {
		l.Error(fmt.Sprintf("task [%s] finished with error ", a.Task.Task.TaskType.String()), zap.Error(err))
	} else {
		if err := a.ProcessRepository.UpdateProcessTask(ctx, a.Task, taskId, a.ProcessId, a.ProfileId); err != nil {
			l.Error("arg.ProcessRepository.UpdateProcessTask", err)
		}
		l.Debug("task finished")
	}

	if err != nil {
		// если транзакция еще не обработалась
		if errors.Is(err, ErrTransactionIsNotReady) {
			return a.Task, nil
		}
		e := err.Error()
		a.Task.Error = &e
		a.Task.Status = v1.ProcessStatus_StatusError

		if err := a.UpdateTask(pctx, a.Task); err != nil {
			return nil, err
		}
	} else {
		a.Task.Error = nil
		if err := a.UpdateTask(pctx, a.Task); err != nil {
			return nil, err
		}
	}

	return task, err
}

func IsPayableTask(t v1.TaskType) bool {

	if config.CFG.Env == config.Local {
		return true
	}

	for _, tt := range PayableTasks {
		if t == tt {
			return true
		}
	}
	return false
}

func NeedPay(before, after *v1.ProcessTask) bool {
	if before.Status == after.Status {
		return false
	}

	if after.Status == v1.ProcessStatus_StatusDone {
		return false
	}

	if !IsPayableTask(after.Task.TaskType) {
		return false
	}

	return true
}

func Marshal(m proto.Message) ([]byte, error) {
	return protojson.MarshalOptions{Multiline: true, UseEnumNumbers: false, UseProtoNames: false, EmitUnpopulated: true}.Marshal(m)
}
