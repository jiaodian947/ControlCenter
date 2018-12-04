package upgrader

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//kernel not found
type RegReplace struct {
	re      *regexp.Regexp
	replace []byte
}

func (rr *RegReplace) Replace(line []byte) []byte {
	line = rr.re.ReplaceAll(line, rr.replace)
	return line
}

type RegRep []*RegReplace

func (rr RegRep) Replace(line []byte) []byte {
	for _, v := range rr {
		line = v.Replace(line)
	}

	return line
}

var kn = []string{
	`pKernel->GetMapBound\(`, `kapi->__GetMapBound(`,
	`pKernel->GetMapHeight\(`, `kapi->__GetMapHeight(`,
	`pKernel->GetMapRegion\(`, `kapi->__GetMapRegion(`,
	`pKernel->GetMapArea\(`, `kapi->__GetMapArea(`,
	`pKernel->GetMapType\(`, `kapi->__GetMapType(`,
	`pKernel->GetWalkType\(`, `kapi->__GetWalkType(`,
	`pKernel->CanWalk\(`, `kapi->__CanWalk(`,
	`pKernel->ObjectCanWalk\(`, `kapi->__ObjectCanWalk(`,
	`pKernel->LineCanWalk\(`, `kapi->__LineCanWalk(`,
	`pKernel->ObjectLineCanWalk\(`, `kapi->__ObjectLineCanWalk(`,
	`pKernel->TraceLineWalk\(`, `kapi->__TraceLineWalk(`,
	`pKernel->GetCollideEnable\(`, `kapi->__GetCollideEnable(`,
	`pKernel->GetApexHeight\(`, `kapi->__GetApexHeight(`,
	`pKernel->GetApexFloor\(`, `kapi->__GetApexFloor(`,
	`pKernel->GetWalkEnable\(`, `kapi->__GetWalkEnable(`,
	`pKernel->GetWalkHeight\(`, `kapi->__GetWalkHeight(`,
	`pKernel->GetWalkWaterExists\(`, `kapi->__GetWalkWaterExists(`,
	`pKernel->GetWalkWaterHeight\(`, `kapi->__GetWalkWaterHeight(`,
	`pKernel->GetFloorCount\(`, `kapi->__GetFloorCount(`,
	`pKernel->GetFloorExists\(`, `kapi->__GetFloorExists(`,
	`pKernel->GetFloorCanMove\(`, `kapi->__GetFloorCanMove(`,
	`pKernel->GetFloorCanStand\(`, `kapi->__GetFloorCanStand(`,
	`pKernel->GetFloorHeight\(`, `kapi->__GetFloorHeight(`,
	`pKernel->GetFloorSpace\(`, `kapi->__GetFloorSpace(`,
	`pKernel->GetFloorHasWall\(`, `kapi->__GetFloorHasWall(`,
	`pKernel->GetWallExists\(`, `kapi->__GetWallExists(`,
	`pKernel->SavePoint\(`, `kapi->__SavePoint(`,
	`pKernel->ResetPoint\(`, `kapi->__ResetPoint(`,
	`pKernel->DeletePoint\(`, `kapi->__DeletePoint(`,
	`pKernel->GetPointCoord\(`, `kapi->__GetPointCoord(`,
	`pKernel->GetScenePointList\(`, `kapi->__GetScenePointList(`,
	`pKernel->LoadObjects\(`, `kapi->__LoadObjects(`,
	`pKernel->SaveObject\(`, `kapi->__SaveObject(`,
	`pKernel->ResetObject\(`, `kapi->__ResetObject(`,
	`pKernel->DeleteObject\(`, `kapi->__DeleteObject(`,
	`pKernel->PreloadConfigTxt\(`, `kapi->__PreloadConfigTxt(`,
	`pKernel->FindConfigProperty\(`, `kapi->__FindConfigProperty(`,
	`pKernel->SetConfigProperty\(`, `kapi->__SetConfigProperty(`,
	`pKernel->GetConfigPropertyList\(`, `kapi->__GetConfigPropertyList(`,
	`pKernel->FindConfigRecord\(`, `kapi->__FindConfigRecord(`,
	`pKernel->GetConfigRecord\(`, `kapi->__GetConfigRecord(`,
	`pKernel->SetConfigRecord\(`, `kapi->__SetConfigRecord(`,
	`pKernel->GetConfigRecordList\(`, `kapi->__GetConfigRecordList(`,
	`pKernel->MoveTo\(`, `kapi->__MoveTo(`,
	`pKernel->Locate\(`, `kapi->__Locate(`,
	`pKernel->Rotate\(`, `kapi->__Rotate(`,
	`pKernel->Motion\(`, `kapi->__Motion(`,
	`pKernel->MotionNoRotate\(`, `kapi->__MotionNoRotate(`,
	`pKernel->Jump\(`, `kapi->__Jump(`,
	`pKernel->JumpTo\(`, `kapi->__JumpTo(`,
	`pKernel->Fly\(`, `kapi->__Fly(`,
	`pKernel->Swim\(`, `kapi->__Swim(`,
	`pKernel->Drift\(`, `kapi->__Drift(`,
	`pKernel->Climb\(`, `kapi->__Climb(`,
	`pKernel->Slide\(`, `kapi->__Slide(`,
	`pKernel->Sink\(`, `kapi->__Sink(`,
	`pKernel->Stop\(`, `kapi->__Stop(`,
	`pKernel->StopWalk\(`, `kapi->__StopWalk(`,
	`pKernel->StopRotate\(`, `kapi->__StopRotate(`,
	`pKernel->CheckMotion\(`, `kapi->__CheckMotion(`,
	`pKernel->CheckJump\(`, `kapi->__CheckJump(`,
	`pKernel->CheckJumpTo\(`, `kapi->__CheckJumpTo(`,
	`pKernel->CheckFly\(`, `kapi->__CheckFly(`,
	`pKernel->CheckSwim\(`, `kapi->__CheckSwim(`,
	`pKernel->CheckDrift\(`, `kapi->__CheckDrift(`,
	`pKernel->CheckClimb\(`, `kapi->__CheckClimb(`,
	`pKernel->CheckSlide\(`, `kapi->__CheckSlide(`,
	`pKernel->CheckSink\(`, `kapi->__CheckSink(`,
	`pKernel->CheckStop\(`, `kapi->__CheckStop(`,
	`pKernel->LinkTo\(`, `kapi->__LinkTo(`,
	`pKernel->LinkMove\(`, `kapi->__LinkMove(`,
	`pKernel->Unlink\(`, `kapi->__Unlink(`,
	`pKernel->GetLinkObject\(`, `kapi->__GetLinkObject(`,
	`pKernel->GetLinkPosition\(`, `kapi->__GetLinkPosition(`,
	`pKernel->GetFloor\(`, `kapi->__GetFloor(`,
	`pKernel->GetDestX\(`, `kapi->__GetDestX(`,
	`pKernel->GetDestY\(`, `kapi->__GetDestY(`,
	`pKernel->GetDestZ\(`, `kapi->__GetDestZ(`,
	`pKernel->GetDestOrient\(`, `kapi->__GetDestOrient(`,
	`pKernel->GetMoveMode\(`, `kapi->__GetMoveMode(`,
	`pKernel->GetMoveSpeed\(`, `kapi->__GetMoveSpeed(`,
	`pKernel->GetRotateSpeed\(`, `kapi->__GetRotateSpeed(`,
	`pKernel->GetJumpSpeed\(`, `kapi->__GetJumpSpeed(`,
	`pKernel->GetGravity\(`, `kapi->__GetGravity(`,
	`pKernel->RequestRecreatePlayer\(`, `kapi->__RequestRecreatePlayer(`,
	`pKernel->ResumePlayer\(`, `kapi->__ResumePlayer(`,
	`pKernel->QueryGift\(`, `kapi->__QueryGift(`,
	`pKernel->GetGift\(`, `kapi->__GetGift(`,
	`pKernel->SellCard\(`, `kapi->__SellCard(`,
	`pKernel->UnsellCard\(`, `kapi->__UnsellCard(`,
	`pKernel->BuyCard\(`, `kapi->__BuyCard(`,
	`pKernel->BuyItem\(`, `kapi->__BuyItem(`,
	`pKernel->BuyItemGive\(`, `kapi->__BuyItemGive(`,
	`pKernel->BuyItem2\(`, `kapi->__BuyItem2(`,
	`pKernel->RequestPoints\(`, `kapi->__RequestPoints(`,
	`pKernel->RequestLimitTime\(`, `kapi->__RequestLimitTime(`,
	`pKernel->RequestAccountInfo\(`, `kapi->__RequestAccountInfo(`,
	`pKernel->RequestPlayTime\(`, `kapi->__RequestPlayTime(`,
	`pKernel->RequestAllItemInfo\(`, `kapi->__RequestAllItemInfo(`,
	`pKernel->RequestAllItemInfo2\(`, `kapi->__RequestAllItemInfo2(`,
	`pKernel->RequestItemInfo\(`, `kapi->__RequestItemInfo(`,
	`pKernel->FindItemNo\(`, `kapi->__FindItemNo(`,
	`pKernel->GetItemCount\(`, `kapi->__GetItemCount(`,
	`pKernel->GetItemNoList\(`, `kapi->__GetItemNoList(`,
	`pKernel->GetItemInfo\(`, `kapi->__GetItemInfo(`,
	`pKernel->GetItemPrice\(`, `kapi->__GetItemPrice(`,
	`pKernel->SaveLog\(`, `kapi->__SaveLog(`,
	`pKernel->ItemLog\(`, `kapi->__ItemLog(`,
	`pKernel->ChatLog\(`, `kapi->__ChatLog(`,
	`pKernel->GmLog\(`, `kapi->__GmLog(`,
	`pKernel->CleanLetter\(`, `kapi->__CleanLetter(`,
	`pKernel->CleanLetterBySerial\(`, `kapi->__CleanLetterBySerial(`,
	`pKernel->SendManageMessage\(`, `kapi->__SendManageMessage(`,
	`pKernel->SendExtraMessage\(`, `kapi->__SendExtraMessage(`,
	`pKernel->SendToGmcc\(`, `kapi->__SendToGmcc(`,
	`pKernel->SendGmccCustom\(`, `kapi->__SendGmccCustom(`,
	`pKernel->FindManageInfo\(`, `kapi->__FindManageInfo(`,
	`pKernel->GetManageInfoName\(`, `kapi->__GetManageInfoName(`,
	`pKernel->GetManageInfoType\(`, `kapi->__GetManageInfoType(`,
	`pKernel->GetManageInfoProperty\(`, `kapi->__GetManageInfoProperty(`,
	`pKernel->GetManageInfoStatus\(`, `kapi->__GetManageInfoStatus(`,
	`pKernel->GetManageInfoVersion\(`, `kapi->__GetManageInfoVersion(`,
	`pKernel->GetManageInfoMaxVersion\(`, `kapi->__GetManageInfoMaxVersion(`,
	`pKernel->UpdateManageBase\(`, `kapi->__UpdateManageBase(`,
	`pKernel->UpdateManageInfo\(`, `kapi->__UpdateManageInfo(`,
	`pKernel->RemoveManageInfo\(`, `kapi->__RemoveManageInfo(`,
	`pKernel->GetManageInfoList\(`, `kapi->__GetManageInfoList(`,
	`pKernel->GetManageInfoListAboveVersion\(`, `kapi->__GetManageInfoListAboveVersion(`,
	`pKernel->RegisteSceneStage\(`, `kapi->__RegisteSceneStage(`,
	`pKernel->GetSceneStage\(`, `kapi->__GetSceneStage(`,
	`pKernel->BindEventProcessor\(`, `kapi->__BindEventProcessor(`,
	`pKernel->GetEventProcessor\(`, `kapi->__GetEventProcessor(`,
	`pKernel->AddEventCallback\(`, `kapi->__AddEventCallback(`,
	`pKernel->RunEvent\(`, `kapi->__RunEvent(`,
	`pKernel->SendEvent\(`, `kapi->__SendEvent(`,
	`pKernel->AddCommandCallback\(`, `kapi->__AddCommandCallback(`,
	`pKernel->RunCommand\(`, `kapi->__RunCommand(`,
	`pKernel->AddCondCallback\(`, `kapi->__AddCondCallback(`,
	`pKernel->RunCondition\(`, `kapi->__RunCondition(`,
	`pKernel->LoadCollideData\(`, `kapi->__LoadCollideData(`,
	`pKernel->TraceLineWalkFromObj\(`, `kapi->__TraceLineWalkFromObj(`,
	`pKernel->GetWalkEnableFromObj\(`, `kapi->__GetWalkEnableFromObj(`,
	`pKernel->GetWalkHeightFromObj\(`, `kapi->__GetWalkHeightFromObj(`,
	`pKernel->GetWalkWaterExistsFromObj\(`, `kapi->__GetWalkWaterExistsFromObj(`,
	`pKernel->GetWalkWaterHeightFromObj\(`, `kapi->__GetWalkWaterHeightFromObj(`,
	`pKernel->GetFloorCountFromObj\(`, `kapi->__GetFloorCountFromObj(`,
	`pKernel->GetFloorExistsFromObj\(`, `kapi->__GetFloorExistsFromObj(`,
	`pKernel->GetFloorCanMoveFromObj\(`, `kapi->__GetFloorCanMoveFromObj(`,
	`pKernel->GetFloorCanStandFromObj\(`, `kapi->__GetFloorCanStandFromObj(`,
	`pKernel->GetFloorHeightFromObj\(`, `kapi->__GetFloorHeightFromObj(`,
	`pKernel->GetFloorSpaceFromObj\(`, `kapi->__GetFloorSpaceFromObj(`,
	`pKernel->GetFloorHasWallFromObj\(`, `kapi->__GetFloorHasWallFromObj(`,
	`pKernel->GetWallExistsFromObj\(`, `kapi->__GetWallExistsFromObj(`,
	`pKernel->SendAuctionMessage\(`, `kapi->__SendAuctionMessage(`,
	`pKernel->SendCrossServerMessage\(`, `kapi->__SendCrossServerMessage(`,
	`pKernel->SendCrossServerRouteMessage\(`, `kapi->__SendCrossServerRouteMessage(`,
	`pKernel->SwitchToBattlefield\(`, `kapi->__SwitchToBattlefield(`,
	`pKernel->QuitBattlefield\(`, `kapi->__QuitBattlefield(`,
	`pKernel->SendReportMessage\(`, `kapi->__SendReportMessage(`,
}
var regkn RegRep

// kernel upgrade
var ku = []string{
	`pKernel->GetResourcePath\(`, `kapi->get_cfg_service()->get_res_path(`,
	`pKernel->GetLogicModule\(`, `kapi->get_logic_module(`,
	`pKernel->AddLogicClass\(`, `kapi->add_logic_class(`,
	`pKernel->AddClassCallback\(`, `kapi->add_class_callback(`,
	`pKernel->AddEventCallback\(`, `kapi->add_event_callback(`,
	`pKernel->AddCommandHook\(`, `kapi->add_command_hook(`,
	`pKernel->AddIntCommandHook\(`, `kapi->add_intcommand_hook(`,
	`pKernel->AddCustomHook\(`, `kapi->add_custom_hook(`,
	`pKernel->AddIntCustomHook\(`, `kapi->add_intcustom_hook(`,
	`pKernel->RunEventCallback\(`, `kapi->run_event_callback(`,
	`pKernel->DeclareLuaExt\(`, `kapi->get_lua_service()->declare_lua_func(`,
	`pKernel->DeclareHeartBeat\(`, `kapi->declare_heartbeat_func(`,
	`pKernel->DeclareCritical\(`, `kapi->declare_attr_hook_func(`,
	`pKernel->DeclareRecHook\(`, `kapi->declare_datagrid_hook_func(`,
	`pKernel->LookupHeartBeat\(`, `kapi->lookup_heartbeat(`,
	`pKernel->LookupCritical\(`, `kapi->lookup_attr_hook(`,
	`pKernel->LookupRecHook\(`, `kapi->lookup_datagrid_hook_func(`,
	`pKernel->Trace\(`, `kapi->log_to_file(`,
	`pKernel->Echo\(`, `kapi->log_to_console(`,
	`pKernel->CheckName\(`, `kapi->check_name(`,
	`pKernel->GetGameId\(`, `kapi->get_cfg_service()->get_game_id(`,
	`pKernel->GetDistrictId\(`, `kapi->get_cfg_service()->get_region_id(`,
	`pKernel->GetServerId\(`, `kapi->get_cfg_service()->get_server_id(`,
	`pKernel->GetMemberId\(`, `kapi->get_cfg_service()->get_scene_server_id(`,
	`pKernel->GetGameObj\(`, `kapi->get_obj(`,
	`pKernel->GetSceneObj\(`, `kapi->get_current_scene_obj(`,
	`pKernel->GetScene\(`, `kapi->get_current_scene_obj_id(`,
	`pKernel->GetSceneId\(`, `kapi->get_current_scene_id(`,
	`pKernel->GetSceneMaxId\(`, `kapi->get_max_normal_scene_id(`,
	`pKernel->GetSceneScript\(`, `kapi->get_scene_script(`,
	`pKernel->GetSceneConfig\(`, `kapi->get_scene_config(`,
	`pKernel->FindSceneId\(`, `kapi->find_scene_id(`,
	`pKernel->GetSceneExists\(`, `kapi->is_scene_exists(`,
	`pKernel->GetSceneOnlineCount\(`, `kapi->get_online_role_number(`,
	`pKernel->GetScenePlayerCount\(`, `kapi->get_all_role_number(`,
	`pKernel->GetScenePlayerList\(`, `kapi->get_role_list(`,
	`pKernel->GetSceneClass\(`, `kapi->get_scene_type(`,
	`pKernel->RequestCloneScene\(`, `kapi->request_clone_scene(`,
	`pKernel->SetCloneSceneDownTime\(`, `kapi->set_clonescene_closetime(`,
	`pKernel->GetCloneSceneCount\(`, `kapi->get_clone_scene_number(`,
	`pKernel->GetCloneSceneList\(`, `kapi->get_clone_scene_list(`,
	`pKernel->GetPrototypeSceneId\(`, `kapi->get_prototype_sceneid(`,
	`pKernel->IsPrototypeScene\(`, `kapi->is_prototype_scene(`,
	`pKernel->SeekRoleUid\(`, `kapi->find_role_guid(`,
	`pKernel->SeekRoleName\(`, `kapi->find_role_name(`,
	`pKernel->GetRoleDeleted\(`, `kapi->is_role_deleted(`,
	`pKernel->CreateObject\(`, `kapi->create_object(`,
	`pKernel->CreateObjectArgs\(`, `kapi->create_object_with_args(`,
	`pKernel->PreloadConfig\(`, `kapi->preload_config_entry(`,
	`pKernel->GetConfigProperty\(`, `kapi->get_config_entry_prop(`,
	`pKernel->LoadConfig\(`, `kapi->load_config_entry(`,
	`pKernel->CreateFromConfig\(`, `kapi->createobj_withconfig_incontainer(`,
	`pKernel->CreateFromConfigArgs\(`, `kapi->createobj_withargs_incontainer(`,
	`pKernel->GetClassIndex\(`, `kapi->get_class_index(`,
	`pKernel->GetClassAttrIndex\(`, `kapi->get_class_attr_index(`,
	`pKernel->Add\(`, `kapi->add_attr(`,
	`pKernel->AddVisible\(`, `kapi->add_visible(`,
	`pKernel->SetVisible\(`, `kapi->set_visible(`,
	`pKernel->SetRealtime\(`, `kapi->set_realtime(`,
	`pKernel->SetSaving\(`, `kapi->set_saving(`,
	`pKernel->Find\(`, `kapi->find_attr(`,
	`pKernel->GetType\(`, `kapi->get_type(`,
	`pKernel->GetAttrVisible\(`, `kapi->get_attr_visible(`,
	`pKernel->GetAttrPublicVisible\(`, `kapi->get_attr_public_visible(`,
	`pKernel->GetAttrRealtime\(`, `kapi->get_attr_realtime(`,
	`pKernel->GetAttrSaving\(`, `kapi->get_attr_saving(`,
	`pKernel->SetAttrHide\(`, `kapi->set_attr_hide(`,
	`pKernel->GetAttrHide\(`, `kapi->get_attr_hide(`,
	`pKernel->GetAttrCount\(`, `kapi->get_attr_count(`,
	`pKernel->GetAttrList\(`, `kapi->get_attr_list(`,
	`pKernel->SetInt\(`, `kapi->set_int(`,
	`pKernel->SetInt64\(`, `kapi->set_int64(`,
	`pKernel->SetFloat\(`, `kapi->set_float(`,
	`pKernel->SetDouble\(`, `kapi->set_double(`,
	`pKernel->SetString\(`, `kapi->set_string(`,
	`pKernel->SetWideStr\(`, `kapi->set_string(`,
	`pKernel->SetObject\(`, `kapi->set_object(`,
	`pKernel->QueryInt\(`, `kapi->get_int(`,
	`pKernel->QueryInt64\(`, `kapi->get_int64(`,
	`pKernel->QueryFloat\(`, `kapi->get_float(`,
	`pKernel->QueryDouble\(`, `kapi->get_double(`,
	`pKernel->QueryString\(`, `kapi->get_string(`,
	`pKernel->QueryWideStr\(`, `kapi->get_string(`,
	`pKernel->QueryObject\(`, `kapi->get_object(`,
	`pKernel->IncInt\(`, `kapi->inc_int(`,
	`pKernel->IncFloat\(`, `kapi->inc_float(`,
	`pKernel->SetAttrFlag\(`, `kapi->set_attr_flag(`,
	`pKernel->ClearAttrFlag\(`, `kapi->clear_attr_flag(`,
	`pKernel->TestAttrFlag\(`, `kapi->test_attr_flag(`,
	`pKernel->AddRecord\(`, `kapi->add_datagrid(`,
	`pKernel->SetRecordKey\(`, `kapi->set_datagrid_key(`,
	`pKernel->SetRecordColType\(`, `kapi->set_datagrid_column_type(`,
	`pKernel->SetRecordVisible\(`, `kapi->set_datagrid_visible(`,
	`pKernel->SetRecordColVisType\(`, `kapi->set_datagrid_column_visible_type(`,
	`pKernel->SetRecordSaving\(`, `kapi->set_datagrid_saving(`,
	`pKernel->GetRecord\(`, `kapi->get_datagrid(`,
	`pKernel->FindRecord\(`, `kapi->find_datagrid(`,
	`pKernel->GetRecordCount\(`, `kapi->get_datagrid_count(`,
	`pKernel->GetRecordList\(`, `kapi->get_datagrid_list(`,
	`pKernel->GetRecordVisible\(`, `kapi->get_datagrid_visible(`,
	`pKernel->GetRecordPublicVisible\(`, `kapi->get_datagrid_public_visible(`,
	`pKernel->GetRecordSaving\(`, `kapi->get_datagrid_saving(`,
	`pKernel->GetRecordCols\(`, `kapi->get_datagrid_cols(`,
	`pKernel->GetRecordRows\(`, `kapi->get_datagrid_rows(`,
	`pKernel->GetRecordMax\(`, `kapi->get_datagrid_max(`,
	`pKernel->GetRecordColType\(`, `kapi->get_datagrid_col_type(`,
	`pKernel->AddRecordRow\(`, `kapi->add_datagrid_row(`,
	`pKernel->AddRecordRowValue\(`, `kapi->add_datagrid_row_value(`,
	`pKernel->RemoveRecordRow\(`, `kapi->remove_datagrid_row(`,
	`pKernel->ClearRecord\(`, `kapi->clear_datagrid(`,
	`pKernel->SetRecordRowValue\(`, `kapi->set_datagrid_row_value(`,
	`pKernel->SetRecordInt\(`, `kapi->set_datagrid_int(`,
	`pKernel->SetRecordInt64\(`, `kapi->set_datagrid_int64(`,
	`pKernel->SetRecordFloat\(`, `kapi->set_datagrid_float(`,
	`pKernel->SetRecordDouble\(`, `kapi->set_datagrid_double(`,
	`pKernel->SetRecordString\(`, `kapi->set_datagrid_string(`,
	`pKernel->SetRecordWideStr\(`, `kapi->set_datagrid_string(`,
	`pKernel->SetRecordObject\(`, `kapi->set_datagrid_object(`,
	`pKernel->QueryRecordRowValue\(`, `kapi->get_datagrid_row_value(`,
	`pKernel->QueryRecordInt\(`, `kapi->get_datagrid_int(`,
	`pKernel->QueryRecordInt64\(`, `kapi->get_datagrid_int64(`,
	`pKernel->QueryRecordFloat\(`, `kapi->get_datagrid_float(`,
	`pKernel->QueryRecordDouble\(`, `kapi->get_datagrid_double(`,
	`pKernel->QueryRecordString\(`, `kapi->get_datagrid_string(`,
	`pKernel->QueryRecordWideStr\(`, `kapi->get_datagrid_string(`,
	`pKernel->QueryRecordObject\(`, `kapi->get_datagrid_object(`,
	`pKernel->FindRecordInt\(`, `kapi->find_datagrid_int(`,
	`pKernel->FindRecordInt64\(`, `kapi->find_datagrid_int64(`,
	`pKernel->FindRecordFloat\(`, `kapi->find_datagrid_float(`,
	`pKernel->FindRecordDouble\(`, `kapi->find_datagrid_double(`,
	`pKernel->FindRecordString\(`, `kapi->find_datagrid_string(`,
	`pKernel->FindRecordWideStr\(`, `kapi->find_datagrid_string(`,
	`pKernel->FindRecordObject\(`, `kapi->find_datagrid_object(`,
	`pKernel->FindRecordStringCI\(`, `kapi->find_datagrid_string_ci(`,
	`pKernel->SetRecordFlag\(`, `kapi->set_datagrid_flag(`,
	`pKernel->ClearRecordFlag\(`, `kapi->clear_datagrid_flag(`,
	`pKernel->TestRecordFlag\(`, `kapi->test_datagrid_flag(`,
	`pKernel->FindData\(`, `kapi->find_volatile(`,
	`pKernel->GetDataCount\(`, `kapi->get_volatile_count(`,
	`pKernel->GetDataList\(`, `kapi->get_volatile_list(`,
	`pKernel->AddData\(`, `kapi->add_volatile(`,
	`pKernel->RemoveData\(`, `kapi->remove_volatile(`,
	`pKernel->GetDataType\(`, `kapi->get_volatile_type(`,
	`pKernel->SetDataInt\(`, `kapi->set_volatile_int(`,
	`pKernel->SetDataInt64\(`, `kapi->set_volatile_int64(`,
	`pKernel->SetDataFloat\(`, `kapi->set_volatile_float(`,
	`pKernel->SetDataDouble\(`, `kapi->set_volatile_double(`,
	`pKernel->SetDataString\(`, `kapi->set_volatile_string(`,
	`pKernel->SetDataWideStr\(`, `kapi->set_volatile_string(`,
	`pKernel->SetDataObject\(`, `kapi->set_volatile_object(`,
	`pKernel->SetDataBinary\(`, `kapi->set_volatile_binary(`,
	`pKernel->QueryDataInt\(`, `kapi->get_volatile_int(`,
	`pKernel->QueryDataInt64\(`, `kapi->get_volatile_int64(`,
	`pKernel->QueryDataFloat\(`, `kapi->get_volatile_float(`,
	`pKernel->QueryDataDouble\(`, `kapi->get_volatile_double(`,
	`pKernel->QueryDataString\(`, `kapi->get_volatile_string(`,
	`pKernel->QueryDataWideStr\(`, `kapi->get_volatile_string(`,
	`pKernel->QueryDataObject\(`, `kapi->get_volatile_object(`,
	`pKernel->QueryDataBinary\(`, `kapi->get_volatile_binary(`,
	`pKernel->GetSwitchLocation\(`, `kapi->get_switch_location(`,
	`pKernel->SwitchLocate\(`, `kapi->switch_to_position(`,
	`pKernel->SwitchBorn\(`, `kapi->switch_to_born_position(`,
	`pKernel->SetSceneBorn\(`, `kapi->set_scene_born_position(`,
	`pKernel->GetSceneBorn\(`, `kapi->get_scene_born_position(`,
	`pKernel->Distance2D\(`, `kapi->distance_2d(`,
	`pKernel->Distance3D\(`, `kapi->distance_3d(`,
	`pKernel->Angle\(`, `kapi->angle(`,
	`pKernel->DotAngle\(`, `kapi->dot_angle(`,
	`pKernel->GetPosiX\(`, `kapi->get_posix(`,
	`pKernel->GetPosiY\(`, `kapi->get_posiy(`,
	`pKernel->GetPosiZ\(`, `kapi->get_posiz(`,
	`pKernel->GetOrient\(`, `kapi->get_orient(`,
	`pKernel->GetLocation\(`, `kapi->get_location(`,
	`pKernel->Create\(`, `kapi->create_obj_by_script(`,
	`pKernel->CreateArgs\(`, `kapi->create_obj_by_args(`,
	`pKernel->CreateContainer\(`, `kapi->create_container(`,
	`pKernel->CreateContainerArgs\(`, `kapi->create_container_by_args(`,
	`pKernel->ExpandContainer\(`, `kapi->expand_container(`,
	`pKernel->ShrinkContainer\(`, `kapi->shrink_container(`,
	`pKernel->CreateClone\(`, `kapi->clone_obj(`,
	`pKernel->CreateTo\(`, `kapi->create_obj_on_position(`,
	`pKernel->CreateToArgs\(`, `kapi->create_obj_on_position_by_args(`,
	`pKernel->Destroy\(`, `kapi->destroy_obj(`,
	`pKernel->DestroySelf\(`, `kapi->destroy_obj_self(`,
	`pKernel->Select\(`, `kapi->select_obj(`,
	`pKernel->SetUnsave\(`, `kapi->set_unsave(`,
	`pKernel->GetUnsave\(`, `kapi->get_unsave(`,
	`pKernel->Place\(`, `kapi->place_obj(`,
	`pKernel->PlacePos\(`, `kapi->place_obj_in_pos(`,
	`pKernel->PlaceTo\(`, `kapi->place_obj_to_position(`,
	`pKernel->Exchange\(`, `kapi->exchange_obj(`,
	`pKernel->ChangePos\(`, `kapi->change_obj_pos(`,
	`pKernel->GetCapacity\(`, `kapi->get_capacity(`,
	`pKernel->GetScript\(`, `kapi->get_script(`,
	`pKernel->GetConfig\(`, `kapi->get_config(`,
	`pKernel->GetIndex\(`, `kapi->get_index(`,
	`pKernel->Type\(`, `kapi->obj_type(`,
	`pKernel->Exists\(`, `kapi->obj_exists(`,
	`pKernel->Parent\(`, `kapi->obj_parent(`,
	`pKernel->GetWeakRefs\(`, `kapi->get_container_refs(`,
	`pKernel->GetWeakBoxList\(`, `kapi->get_container_list(`,
	`pKernel->GetItem\(`, `kapi->get_item(`,
	`pKernel->GetChild\(`, `kapi->get_child(`,
	`pKernel->GetFirst\(`, `kapi->get_first(`,
	`pKernel->GetNext\(`, `kapi->get_next(`,
	`pKernel->GetChildCount\(`, `kapi->get_child_count(`,
	`pKernel->GetChildList\(`, `kapi->get_child_list(`,
	`pKernel->ClearChild\(`, `kapi->clear_child(`,
	`pKernel->FindChild\(`, `kapi->find_child(`,
	`pKernel->FindChildMore\(`, `kapi->find_child_more(`,
	`pKernel->FindChildByConfig\(`, `kapi->find_child_by_config(`,
	`pKernel->FindChildMoreByConfig\(`, `kapi->find_child_more_by_config(`,
	`pKernel->AddWeakChild\(`, `kapi->add_child_to_container(`,
	`pKernel->RemoveWeakChild\(`, `kapi->remove_child_from_container(`,
	`pKernel->ClearWeakChild\(`, `kapi->clear_child_from_container(`,
	`pKernel->GetAroundList\(`, `kapi->get_around_obj_list(`,
	`pKernel->GetPointAroundList\(`, `kapi->get_around_point_list(`,
	`pKernel->GetAroundList3D\(`, `kapi->get_around_obj_list_3d(`,
	`pKernel->GetPointAroundList3D\(`, `kapi->get_around_point_list_3d(`,
	`pKernel->TraceObjectList\(`, `kapi->segment_line_trace_list(`,
	`pKernel->NewGroupId\(`, `kapi->generate_new_group_id(`,
	`pKernel->GroupFirst\(`, `kapi->get_group_first(`,
	`pKernel->GroupNext\(`, `kapi->get_group_next(`,
	`pKernel->GroupChildList\(`, `kapi->get_group_child_obj_list(`,
	`pKernel->GroupFind\(`, `kapi->search_obj_by_name_in_group(`,
	`pKernel->GroupFindMore\(`, `kapi->search_all_obj_by_name_in_group(`,
	`pKernel->ClearGroup\(`, `kapi->clear_non_role_obj_in_group(`,
	`pKernel->AddHeartBeat\(`, `kapi->add_heart_beat_to_obj(`,
	`pKernel->AddCountBeat\(`, `kapi->add_count_heart_beat_to_obj(`,
	`pKernel->RemoveHeartBeat\(`, `kapi->remove_heart_beat_from_obj(`,
	`pKernel->ClearHeartBeat\(`, `kapi->clear_heart_beat_from_obj(`,
	`pKernel->FindHeartBeat\(`, `kapi->search_heart_beat_by_name(`,
	`pKernel->GetBeatTime\(`, `kapi->get_obj_hear_beat_time_by_name(`,
	`pKernel->SetBeatCount\(`, `kapi->set_obj_heart_beat_count(`,
	`pKernel->GetBeatCount\(`, `kapi->get_obj_heart_beat_count(`,
	`pKernel->FindCritical\(`, `kapi->attr_has_hook_func(`,
	`pKernel->AddCritical\(`, `kapi->add_attr_hook_func(`,
	`pKernel->RemoveCritical\(`, `kapi->remove_attr_hook_func_all(`,
	`pKernel->RemoveCriticalFunc\(`, `kapi->remove_attr_hook_func(`,
	`pKernel->FindRecHook\(`, `kapi->find_datagrid_hook(`,
	`pKernel->AddRecHook\(`, `kapi->add_datagrid_hook(`,
	`pKernel->RemoveRecHook\(`, `kapi->remove_datagrid_hook_all(`,
	`pKernel->RemoveRecHookFunc\(`, `kapi->remove_datagrid_hook_func(`,
	`pKernel->AddViewport\(`, `kapi->add_monitor(`,
	`pKernel->RemoveViewport\(`, `kapi->remove_monitor(`,
	`pKernel->FindViewport\(`, `kapi->find_monitor(`,
	`pKernel->GetViewportContainer\(`, `kapi->get_monitor_container(`,
	`pKernel->ClearViewport\(`, `kapi->clear_monitor(`,
	`pKernel->GetViewers\(`, `kapi->get_monitors(`,
	`pKernel->CloseViewers\(`, `kapi->close_monitors(`,
	`pKernel->Speech\(`, `kapi->talk_to_around_players(`,
	`pKernel->SetLandPoint\(`, `kapi->set_land_point(`,
	`pKernel->SetLandPosi\(`, `kapi->set_land_posi(`,
	`pKernel->RequestAccountRole\(`, `kapi->request_role_account(`,
	`pKernel->RequestRoleInfo\(`, `kapi->request_role_preview_info(`,
	`pKernel->RequestEditPlayer\(`, `kapi->request_edit_role(`,
	`pKernel->BreakPlayer\(`, `kapi->break_player_by_guid(`,
	`pKernel->BreakByAccount\(`, `kapi->break_player_by_account(`,
	`pKernel->BreakByName\(`, `kapi->break_player_by_name(`,
	`pKernel->BlockPlayer\(`, `kapi->block_player_by_name(`,
	`pKernel->SetOfflineTime\(`, `kapi->set_player_offline_live_time(`,
	`pKernel->PlayerLeaveGame\(`, `kapi->player_request_leave_game(`,
	`pKernel->FindPlayer\(`, `kapi->search_player_by_name(`,
	`pKernel->GetPlayerScene\(`, `kapi->get_player_scene_id(`,
	`pKernel->GetPlayerCount\(`, `kapi->get_player_number(`,
	`pKernel->GetOnlineCount\(`, `kapi->get_online_player_number(`,
	`pKernel->GetOfflineCount\(`, `kapi->get_offline_player_number(`,
	`pKernel->ChannelAdd\(`, `kapi->add_player_to_channel(`,
	`pKernel->ChannelRemove\(`, `kapi->remove_player_from_channel(`,
	`pKernel->SysInfo\(`, `kapi->send_sysinfo_by_id(`,
	`pKernel->SysInfoByName\(`, `kapi->send_sysinfo_by_name(`,
	`pKernel->SysInfoByKen\(`, `kapi->send_sysinfo_to_players_around(`,
	`pKernel->SysInfoByScene\(`, `kapi->send_sysinfo_to_players_in_scene(`,
	`pKernel->SysInfoByGroup\(`, `kapi->send_sysinfo_to_players_in_group(`,
	`pKernel->SysInfoByWorld\(`, `kapi->send_sysinfo_to_players_in_world(`,
	`pKernel->SysInfoByChannel\(`, `kapi->send_sysinfo_to_players_in_channel(`,
	`pKernel->SysInfoBroadcast\(`, `kapi->send_sysinfo_to_all_players(`,
	`pKernel->Custom\(`, `kapi->send_custom_msg_by_id(`,
	`pKernel->Custom2\(`, `kapi->send_custom_msg_by_id_2(`,
	`pKernel->CustomByName\(`, `kapi->send_custom_msg_by_name(`,
	`pKernel->CustomByName2\(`, `kapi->send_custom_msg_by_name_2(`,
	`pKernel->CustomByKen\(`, `kapi->send_custom_msg_to_around_players(`,
	`pKernel->CustomByScene\(`, `kapi->send_custom_msg_to_players_in_scene(`,
	`pKernel->CustomByGroup\(`, `kapi->send_custom_msg_to_players_in_group(`,
	`pKernel->CustomByWorld\(`, `kapi->send_custom_msg_to_players_in_world(`,
	`pKernel->CustomByChannel\(`, `kapi->send_custom_msg_to_players_in_channel(`,
	`pKernel->Command\(`, `kapi->send_command_msg_to_obj_by_id(`,
	`pKernel->CommandByName\(`, `kapi->send_command_msg_to_obj_by_name(`,
	`pKernel->CommandByKen\(`, `kapi->send_command_msg_to_objs_around(`,
	`pKernel->CommandByScene\(`, `kapi->send_command_msg_to_objs_in_scene(`,
	`pKernel->CommandByGroup\(`, `kapi->send_command_msg_to_objs_in_group(`,
	`pKernel->CommandBySceneGroup\(`, `kapi->send_command_msg_to_objs_in_scene_group(`,
	`pKernel->CommandByWorld\(`, `kapi->send_command_msg_to_objs_in_world(`,
	`pKernel->CommandToScene\(`, `kapi->send_command_msg_to_scene(`,
	`pKernel->CommandToAllScene\(`, `kapi->send_command_msg_to_all_scenes(`,
	`pKernel->SavePlayerData\(`, `kapi->save_player_data(`,
	`pKernel->FindChunk\(`, `kapi->search_scene_chunk_by_name(`,
	`pKernel->SaveChunk\(`, `kapi->save_scene_chunk(`,
	`pKernel->LoadChunk\(`, `kapi->load_scene_chunk(`,
	`pKernel->DeleteChunk\(`, `kapi->delete_scene_chunk(`,
	`pKernel->GetChunkNameList\(`, `kapi->get_scene_chunk_name_list(`,
	`pKernel->GetChunkObjectClass\(`, `kapi->get_scene_chunk_object_class(`,
	`pKernel->GetChunkObjectScript\(`, `kapi->get_scene_chunk_object_script(`,
	`pKernel->GetChunkObjectConfig\(`, `kapi->get_scene_chunk_object_config(`,
	`pKernel->ReleaseAllChunk\(`, `kapi->release_all_scene_chunk(`,
	`pKernel->CustomLog\(`, `kapi->save_custom_log(`,
	`pKernel->CustomLogWithRole\(`, `kapi->save_custom_log_with_role(`,
	`pKernel->SendLetter\(`, `kapi->send_letter_to_player(`,
	`pKernel->SystemLetter\(`, `kapi->system_letter_to_player(`,
	`pKernel->SystemLetterByAccount\(`, `kapi->system_letter_to_player_by_account(`,
	`pKernel->RecvLetter\(`, `kapi->recv_letter_and_delete(`,
	`pKernel->RecvLetterBySerial\(`, `kapi->recv_letter_delete_by_letter_id(`,
	`pKernel->LookLetter\(`, `kapi->look_letter(`,
	`pKernel->QueryLetter\(`, `kapi->query_letter(`,
	`pKernel->BackLetterBySerial\(`, `kapi->back_letter_by_letter_id(`,
	`pKernel->LuaLoadScript\(`, `kapi->get_lua_service()->load_lua_script(`,
	`pKernel->LuaFindScript\(`, `kapi->get_lua_service()->search_lua_script(`,
	`pKernel->LuaFindFunc\(`, `kapi->get_lua_service()->search_lua_func(`,
	`pKernel->LuaRunFunc\(`, `kapi->get_lua_service()->run_lua_func(`,
	`pKernel->LuaErrorHandler\(`, `kapi->get_lua_service()->do_lua_error_handler(`,
	`pKernel->LuaGetArgCount\(`, `kapi->get_lua_service()->get_lua_arg_count(`,
	`pKernel->LuaIsInt\(`, `kapi->get_lua_service()->lua_is_int(`,
	`pKernel->LuaIsInt64\(`, `kapi->get_lua_service()->lua_is_int64(`,
	`pKernel->LuaIsFloat\(`, `kapi->get_lua_service()->lua_is_float(`,
	`pKernel->LuaIsDouble\(`, `kapi->get_lua_service()->lua_is_double(`,
	`pKernel->LuaIsString\(`, `kapi->get_lua_service()->lua_is_string(`,
	`pKernel->LuaIsWideStr\(`, `kapi->get_lua_service()->lua_is_string(`,
	`pKernel->LuaIsObject\(`, `kapi->get_lua_service()->lua_is_object(`,
	`pKernel->LuaToInt\(`, `kapi->get_lua_service()->lua_get_int(`,
	`pKernel->LuaToInt64\(`, `kapi->get_lua_service()->lua_get_int64(`,
	`pKernel->LuaToFloat\(`, `kapi->get_lua_service()->lua_get_float(`,
	`pKernel->LuaToDouble\(`, `kapi->get_lua_service()->lua_get_double(`,
	`pKernel->LuaToString\(`, `kapi->get_lua_service()->lua_get_string(`,
	`pKernel->LuaToWideStr\(`, `kapi->get_lua_service()->lua_get_string(`,
	`pKernel->LuaToObject\(`, `kapi->get_lua_service()->lua_get_object(`,
	`pKernel->LuaPushBool\(`, `kapi->get_lua_service()->lua_push_bool(`,
	`pKernel->LuaPushNumber\(`, `kapi->get_lua_service()->lua_push_number(`,
	`pKernel->LuaPushInt\(`, `kapi->get_lua_service()->lua_push_int(`,
	`pKernel->LuaPushInt64\(`, `kapi->get_lua_service()->lua_push_int64(`,
	`pKernel->LuaPushFloat\(`, `kapi->get_lua_service()->lua_push_float(`,
	`pKernel->LuaPushDouble\(`, `kapi->get_lua_service()->lua_push_double(`,
	`pKernel->LuaPushString\(`, `kapi->get_lua_service()->lua_push_string(`,
	`pKernel->LuaPushWideStr\(`, `kapi->get_lua_service()->lua_push_string(`,
	`pKernel->LuaPushObject\(`, `kapi->get_lua_service()->lua_push_object(`,
	`pKernel->LuaPushTable\(`, `kapi->get_lua_service()->lua_push_table(`,
	`pKernel->FindPubSpace\(`, `kapi->is_global_data_exist(`,
	`pKernel->GetPubSpaceCount\(`, `kapi->get_global_data_number(`,
	`pKernel->GetPubSpaceList\(`, `kapi->get_global_data_name_list(`,
	`pKernel->GetPubSpace\(`, `kapi->get_global_data(`,
	`pKernel->SendPublicMessage\(`, `kapi->send_msg_to_global_data_server(`,
	`pKernel->SendChargeMessage\(`, `kapi->send_msg_to_bill_server(`,
	`pKernel->SendChargeSafeMessage\(`, `kapi->send_safe_msg_to_bill_server(`,
	`pKernel->SendReportMessage\(`, `kapi->send_report_message(`,
}
var regku RegRep

//game object upgrade
var gu = []string{
	`->GetClassType\(`, `->get_obj_type(`,
	`->GetObjectId\(`, `->get_object_id(`,
	`->GetScript\(`, `->get_script(`,
	`->GetConfig\(`, `->get_config(`,
	`->GetName\(`, `->get_name(`,
	`->GetGroupId\(`, `->get_group_id(`,
	`->GetIndex\(`, `->get_index_in_container(`,
	`->GetParent\(`, `->get_parent_obj(`,
	`->GetCapacity\(`, `->get_container_capacity(`,
	`->GetChildCount\(`, `->get_child_obj_number(`,
	`->GetChildByIndex\(`, `->get_child_obj_by_index(`,
	`->GetChild\(`, `->get_child_obj_by_name(`,
	`->GetWeakRefs\(`, `->get_container_refs(`,
	`->GetPosiX\(`, `->__GetPosiX(`,
	`->GetPosiY\(`, `->__GetPosiY(`,
	`->GetPosiZ\(`, `->__GetPosiZ(`,
	`->GetOrient\(`, `->__GetOrient(`,
	`->FindAttr\(`, `->is_attr_exist(`,
	`->GetAttrType\(`, `->get_attr_type(`,
	`->GetAttrCount\(`, `->get_attr_count(`,
	`->GetAttrList\(`, `->get_attr_name_list(`,
	`->SetInt\(`, `->set_int(`,
	`->SetInt64\(`, `->set_int64(`,
	`->SetFloat\(`, `->set_float(`,
	`->SetDouble\(`, `->set_double(`,
	`->SetString\(`, `->set_string(`,
	`->SetWideStr\(`, `->set_string(`,
	`->SetObject\(`, `->set_object(`,
	`->QueryInt\(`, `->get_int(`,
	`->QueryInt64\(`, `->get_int64(`,
	`->QueryFloat\(`, `->get_float(`,
	`->QueryDouble\(`, `->get_double(`,
	`->QueryString\(`, `->get_string(`,
	`->QueryWideStr\(`, `->get_string(`,
	`->QueryObject\(`, `->get_object(`,
	`->GetAttrIndex\(`, `->get_attr_index(`,
	`->SetIntByIndex\(`, `->set_int(`,
	`->SetInt64ByIndex\(`, `->set_int64(`,
	`->SetFloatByIndex\(`, `->set_float(`,
	`->SetDoubleByIndex\(`, `->set_double(`,
	`->SetStringByIndex\(`, `->set_string(`,
	`->SetWideStrByIndex\(`, `->set_string(`,
	`->SetObjectByIndex\(`, `->set_object(`,
	`->QueryIntByIndex\(`, `->get_int(`,
	`->QueryInt64ByIndex\(`, `->get_int64(`,
	`->QueryFloatByIndex\(`, `->get_float(`,
	`->QueryDoubleByIndex\(`, `->get_double(`,
	`->QueryStringByIndex\(`, `->get_string(`,
	`->QueryWideStrByIndex\(`, `->get_string(`,
	`->QueryObjectByIndex\(`, `->get_object(`,
	`->GetRecordCount\(`, `->get_datagrid_count(`,
	`->GetRecordList\(`, `->get_datagrid_name_list(`,
	`->GetRecordByIndex\(`, `->get_datagrid_by_index(`,
	`->GetRecord\(`, `->get_datagrid(`,
	`->GetRecordIndex\(`, `->get_datagrid_index(`,
	`->FindData\(`, `->volatile_exist(`,
	`->GetDataCount\(`, `->get_volatile_count(`,
	`->GetDataList\(`, `->get_volatile_name_list(`,
	`->AddDataInt\(`, `->add_volatile_int(`,
	`->AddDataInt64\(`, `->add_volatile_int64(`,
	`->AddDataFloat\(`, `->add_volatile_float(`,
	`->AddDataDouble\(`, `->add_volatile_double(`,
	`->AddDataString\(`, `->add_volatile_string(`,
	`->AddDataWideStr\(`, `->add_volatile_string(`,
	`->AddDataObject\(`, `->add_volatile_object(`,
	`->AddDataBinary\(`, `->add_volatile_binary(`,
	`->RemoveData\(`, `->remove_volatile(`,
	`->GetDataType\(`, `->get_volatile_type(`,
	`->SetDataInt\(`, `->set_volatile_int(`,
	`->SetDataInt64\(`, `->set_volatile_int64(`,
	`->SetDataFloat\(`, `->set_volatile_float(`,
	`->SetDataDouble\(`, `->set_volatile_double(`,
	`->SetDataString\(`, `->set_volatile_string(`,
	`->SetDataWideStr\(`, `->set_volatile_string(`,
	`->SetDataObject\(`, `->set_volatile_object(`,
	`->SetDataBinary\(`, `->set_volatile_binary(`,
	`->QueryDataInt\(`, `->get_volatile_int(`,
	`->QueryDataInt64\(`, `->get_volatile_int64(`,
	`->QueryDataFloat\(`, `->get_volatile_float(`,
	`->QueryDataDouble\(`, `->get_volatile_double(`,
	`->QueryDataString\(`, `->get_volatile_string(`,
	`->QueryDataWideStr\(`, `->get_volatile_string(`,
	`->QueryDataObject\(`, `->get_volatile_object(`,
	`->QueryDataBinary\(`, `->get_volatile_binary(`,
	`->FindViewport\(`, `->search_monitor(`,
	`->GetMoveMode\(`, `->__GetMoveMode(`,
}
var reggu RegRep

// callback rename
var cbrename = []string{
	`"OnCreateClass"`, `"on_create_cls"`,
	`"OnCreateRole"`, `"on_create_role"`,
	`"OnCreate"`, `"on_create"`,
	`"OnCreateArgs"`, `"__OnCreateArgs"`,
	`"OnDestroy"`, `"on_destroy"`,
	`"OnSelect"`, `"on_select"`,
	`"OnSpring"`, `"on_spring"`,
	`"OnEndSpring"`, `"on_end_spring"`,
	`"OnEntry"`, `"on_entry"`,
	`"OnLeave"`, `"on_leave"`,
	`"OnNoAdd"`, `"on_noadd"`,
	`"OnAdd"`, `"on_add"`,
	`"OnAfterAdd"`, `"on_after_add"`,
	`"OnNoRemove"`, `"on_noremove"`,
	`"OnBeforeRemove"`, `"on_before_remove"`,
	`"OnRemove"`, `"on_remove"`,
	`"OnChange"`, `"on_change"`,
	`"OnReady"`, `"on_ready"`,
	`"OnCommand"`, `"on_cmd"`,
	`"OnLoad"`, `"on_load"`,
	`"OnRecover"`, `"on_recover"`,
	`"OnStore"`, `"on_store"`,
	`"OnSysInfo"`, `"on_sysinfo"`,
	`"OnSpeech"`, `"on_speech"`,
	`"OnCustom"`, `"on_custom"`,
	`"OnDoSelect"`, `"on_do_select"`,
	`"OnMotion"`, `"on_motion"`,
	`"OnRotate"`, `"on_rotate"`,
	`"OnRequestMove"`, `"on_request_move"`,
	`"OnDisconnect"`, `"on_disconnect"`,
	`"OnContinue"`, `"on_continue"`,
	`"OnEntryGame"`, `"on_entry_game"`,
	`"OnBeforeEntryScene"`, `"on_before_entry_game"`,
	`"OnAfterEntryScene"`, `"on_after_entry_game"`,
	`"OnLeaveScene"`, `"on_leave_scene"`,
	`"OnGetWorldInfo"`, `"on_get_world_info"`,
	`"OnSendLetter"`, `"on_send_letter"`,
	`"OnQueryLetter"`, `"on_query_letter"`,
	`"OnRecvLetter"`, `"on_recv_letter"`,
	`"OnRecvLetterFail"`, `"on_recv_letter_fail"`,
	`"OnLookLetter"`, `"on_look_letter"`,
	`"OnCleanLetter"`, `"on_clean_letter"`,
	`"OnBackLetter"`, `"on_back_letter"`,
	`"OnChargeNotify"`, `"on_bill_srv_notify"`,
	`"OnMapChanged"`, `"on_map_changed"`,
	`"OnMapTypeChanged"`, `"on_map_type_changed"`,
	`"OnEntryVisual"`, `"on_entry_visual"`,
	`"OnLeaveVisual"`, `"on_leave_visual"`,
	`"OnGmccCustom"`, `"on_gmcc_custom"`,
	`"OnSystemStatInfo"`, `"on_sys_statinfo"`,
	`"OnChargeMessage"`, `"on_bill_srv_msg"`,
	`"OnManageMessage"`, `"on_mng_sys_msg"`,
	`"OnExtraMessage"`, `"on_extra_srv_msg"`,
	`"OnPublicMessage"`, `"on_share_srv_msg"`,
	`"OnPublicComplete"`, `"on_share_srv_complete"`,
	`"OnServerClose"`, `"on_srv_close"`,
	`"OnCloneScene"`, `"on_clone_scene"`,
	`"OnCloneReset"`, `"on_clone_reset"`,
	`"OnEditPlayer"`, `"on_edit_player"`,
	`"OnRecreatePlayer"`, `"on_recreate_player"`,
	`"OnGetAccountRole"`, `"on_get_acct_role"`,
	`"OnGetRoleInfo"`, `"on_get_role_info"`,
}
var regcbrename RegRep

var varlistupgrade = []string{
	`AddBool`, `add_bool`,
	`AddInt`, `add_int`,
	`AddInt64`, `add_int64`,
	`AddFloat`, `add_float`,
	`AddDouble`, `add_double`,
	`AddString`, `add_string`,
	`AddWideStr`, `add_string`,
	`AddObject`, `add_object`,
	`AddPointer`, `add_pointer`,
	`AddUserData`, `add_userdata`,
	`AddRawUserData`, `add_raw_userdata`,
	`BoolVal`, `get_bool`,
	`IntVal`, `get_int`,
	`Int64Val`, `get_int64`,
	`FloatVal`, `get_float`,
	`DoubleVal`, `get_double`,
	`StringVal`, `get_string`,
	`WideStrVal`, `get_string`,
	`ObjectVal`, `get_object`,
	`PointerVal`, `get_pointer`,
	`UserDataVal`, `get_userdata`,
	`RawUserDataVal`, `get_raw_userdata`,
}
var regvarlistupgrade RegRep

// class macro rename
var classRename = []string{
	`IKernel`, `KAPI`,
	`ICore`, `CAPI`,
	`IPubKernel`, `GAPI`,
	`IRecord`, `IDataGrid`,
	`IGameObj`, `IBaseObject`,
	`PERSISTID`, `ObjId`,
	`IVarList`, `IArrayList`,
	`TVarList`, `TArrayList`,
	`CVarList`, `ArrayList`,
	`IVar`, `IAny`,
	`TVar`, `TAny`,
	`CVar`, `Any`,
	`TFastStr`, `TStr`,
	`return_string`, `AutoStr`,
	`return_wstring`, `AutoWstr`,
	`fast_string`, `AutoStr`,
	`fast_wstring`, `AutoWstr`,
	`CLoadArchive`, `SerializeReader`,
	`CStoreArchive`, `SerializeWriter`,
	`wchar_t\s*\*`, `char*`,
	`wchar_t\s`, `char `,
}
var regclassrename RegRep

//character rename
var rename = []string{
	`(?i:IKernel\* pKernel)`, `KAPI* kapi`,
	`(?i:KAPI* pKernel)`, `KAPI* kapi`,
	`(?i:pKernel)`, `kapi`,
	`\sInit\(`, ` init(`,
	`\sShut\(`, ` shut(`,
	`\sBeforeLaunch\(`, ` before_launch(`,
	`::Init\(`, `::init(`,
	`::Shut\(`, `::shut(`,
	`::BeforeLaunch\(`, `::before_launch(`,
}
var regrename RegRep

var redefine = []string{
	`FX_MODULE_CORE_VERSION`, `RUNTIME_VERSION`,
	`FX_DLL_EXPORT`, `DLL_EXPORT`,
	`Assert`, `ASSERT`,
	`FxModule_GetType`, `get_module_type`,
	`FxModule_Init`, `logic_module_init`,
	`FxGameLogic_GetVersion`, `get_module_version`,
	`FxGameLogic_GetCreator`, `get_logic_creator`,
	`FxGameLogic_GetModuleCreator`, `get_logic_module_creator`,
	`DeclareLuaExt`, `define_lua_dispatch_func`,
	`DeclareHeartBeat`, `define_hearbeat_dispatch_func`,
	`DeclareCritical`, `define_property_dispatch_func`,
	`DeclareRecHook`, `define_table_dispatch_func`,
	`FX_SYSTEM_WINDOWS`, `_PLATFORM_WINDOWS_`,
	`FX_SYSTEM_32BIT`, `_PLATFORM_X32_`,
	`FX_SYSTEM_64BIT`, `_PLATFORM_X64_`,
	`FX_SYSTEM_LINUX`, `_PLATFORM_LINUX_`,
	`SafeSprintf`, `SAFE_SPRINTF`,
	`SafeSprintList`, `SAFE_SPRINT_LIST`,
	`SafeSwprintf`, `SAFE_SWPRINTF`,
	`DECL_LUA_EXT`, `DECL_LUA_FUNC`,
	`DECL_HEARTBEAT`, `DECL_HEARTBEAT_FUNC`,
	`DECL_CRITICAL`, `DECL_PROPERTY_HOOK_FUNC`,
	`DECL_RECHOOK`, `DECL_TABLE_HOOK_FUNC`,
	`TYPE_SCENE`, `E_OBJ_SCENE`,
	`TYPE_PLAYER`, `E_OBJ_ROLE`,
	`TYPE_NPC`, `E_OBJ_NPC`,
	`TYPE_ITEM`, `E_OBJ_ITEM`,
	`TYPE_HELPER`, `E_OBJ_AIDE`,
	`TYPE_WEAKBOX`, `E_OBJ_CONTAINER`,
	`SHAPE_CYLINDER`, `E_BV_CYLINDER`,
	`SHAPE_SPHERE`, `E_BV_SPHERE`,
	`SHAPE_POLYGON`, `E_BV_POLYGON`,
	`STORE_EXIT`, `E_STORE_EXIT`,
	`STORE_TIMING`, `E_STORE_TIMING`,
	`STORE_SWITCH`, `E_STORE_SWITCH`,
	`STORE_MANUAL`, `E_STORE_MANUAL`,
	`STORE_EDIT`, `E_STORE_EDIT`,
	`STORE_RECREATE`, `E_STORE_RECREATE`,
	`RECOP_INIT`, `E_DATAGRID_INIT`,
	`RECOP_ADD_ROW`, `E_DATAGRID_ADD_ROW`,
	`RECOP_REMOVE_ROW`, `E_DATAGRID_REMOVE_ROW`,
	`RECOP_CLEAR_ROW`, `E_DATAGRID_CLEAR_ROW`,
	`RECOP_GRID_CHANGE`, `E_DATAGRID_GRID_CHANGE`,
	`RECOP_SET_ROW`, `E_DATAGRID_SET_ROW`,
}
var regdefine RegRep

func makere(kv []string) RegRep {
	size := len(kv)
	ret := make(RegRep, 0, size/2)
	for i := 0; i < size; {
		re := &RegReplace{}
		re.re = regexp.MustCompile(kv[i])
		re.replace = []byte(kv[i+1])
		ret = append(ret, re)
		i += 2
	}
	return ret
}

// find path change
var repath = []string{
	`(?i:#include\s+"share/(?P<file>[\w\.]+)")`, `#include <${file}>`,
	`(?i:#include\s+"utility/(?P<file>[\w\.]+)")`, `#include <${file}>`,
	`(?i:#include\s+"define/(?P<file>[\w\.]+)")`, `#include <${file}>`,
	`(?i:#include\s+"public/(?P<file>[\w\.]+)")`, `#include <${file}>`,
	`(?i:#include\s+"utils/(?P<file>[\w\.]+)")`, `#include <${file}>`,
	`(?i:#include\s+"system/(?P<file>[\w\.]+)")`, `#include <${file}>`,
	`(?i:#include\s+"server/(?P<file>[\w\.]+)")`, `#include <${file}>`,
	`(?i:#include\s+"SDK/(?P<file>[\w\.]+)")`, `#include <${file}>`,
}

var regrepath RegRep

// file rename
var filerename = []string{
	`Macros\.h`, `macros_def.h`,
	`ShareMem\.h`, `share_mem.h`,
	`LogFile\.h`, `log_file.h`,
	`IEntity\.h`, `api_base_entity.h`,
	`IEntCreator\.h`, `api_entity_creator.h`,
	`EntManager\.h`, `entity_mng.h`,
	`IEntInfo\.h`, `api_entity_info.h`,
	`EntInfo\.h`, `entity_info.h`,
	`EntInfoList\.h`, `entity_info_mng.h`,
	`EntFactory\.h`, `entity_factory.h`,
	`IEntScript\.h`, `api_entity_script.h`,
	`EntScript\.h`, `entity_script.h`,
	`ICore\.h`, `capi.h`,
	`IInterface\.h`, `api_base_interface.h`,
	`IIntCreator\.h`, `api_base_interface_creator.h`,
	`IntManager\.h`, `interface_mng.h`,
	`ILogic\.h`, `api_base_logic.h`,
	`ILogicInfo\.h`, `api_base_logic_info.h`,
	`LogicInfo\.h`, `logic_info.h`,
	`LogicInfoList\.h`, `logic_info_mng.h`,
	`ILogicCreator\.h`, `api_base_logic_creator.h`,
	`LogicManager\.h`, `logic_mng.h`,
	`ILogicLoader\.h`, `api_logic_mod_loader.h`,
	`"Module\.h"`, `"engine_module_def.h"`,
	`[\\/]Module\.h`, `/engine_module_def.h`,
	`LogicDll\.h`, `logic_module_def.h`,
	`IFuncCreator\.h`, `api_func_creator.h`,
	`FuncManager\.h`, `function_mng.h`,
	`VarTraits\.h`, `var_traits.h`,
	`VarGetter\.h`, `var_getter.h`,
	`VarSetter\.h`, `var_setter.h`,
	`CoreMem\.h`, `core_mem_mng.h`,
	`CoreConfig\.h`, `core_cfg.h`,
	`DllManager\.h`, `module_mng.h`,
	`IConsole\.h`, `api_console.h`,
	`LeakChecker\.h`, `memory_leak_report.h`,
	`MemPool\.h`, `smallblock_pool.h`,
	`Portable\.h`, `crossplatform.h`,
	`Routine\.h`, `function_misc.h`,
	`CharTraits\.h`, `char_traits.h`,
	`Converts\.h`, `string_convert.h`,
	`ServerConst\.h`, `server_const.h`,
	`OuterMsg\.h`, `server_client_base_msg.h`,
	`KnlConst\.h`, `engine_const.h`,
	`PoolArray\.h`, `mng_vector.h`,
	`PoolString\.h`, `mng_string.h`,
	`PoolWideStr\.h`, `mng_wstring.h`,
	`FastStr\.h`, `t_string.h`,
	`ArrayPod\.h`, `t_vector.h`,
	`StringPod\.h`, `hashmap_string.h`,
	`StringMap\.h`, `hashmap_string.h`,
	`PodTraits\.h`, `base_type_traits.h`,
	`PodHashMap\.h`, `t_hashmap.h`,
	`LoadArchive\.h`, `serialize_read.h`,
	`StoreArchive\.h`, `serialize_write.h`,
	`CoreLog\.h`, `log_def.h`,
	`PersistId\.h`, `obj_guid.h`,
	`Location\.h`, `t_position.h`,
	`NameUid\.h`, `name_map_uid.h`,
	`NameFilter\.h`, `name_filter.h`,
	`IVar\.h`, `api_any.h`,
	`Var\.h`, `ls_var.h`,
	`VarType\.h`, `any_type.h`,
	`IVarList\.h`, `api_any_list.h`,
	`VarList\.h`, `any_list.h`,
	`IRecord\.h`, `api_datagrid.h`,
	`IPubData\.h`, `api_global_data.h`,
	`PubData\.h`, `global_data.h`,
	`PublicRecord\.h`, `global_data_table.h`,
	`ILogicModule\.h`, `api_logic_module.h`,
	`IModuleCreator\.h`, `api_logic_module_creator.h`,
	`IKernel\.h`, `kapi.h`,
	`IGameObj\.h`, `api_base_object.h`,
	`LogicCreator\.h`, `logic_class_creator.h`,
	`LogicCreatorSet\.h`, `logic_class_creator_mgr.h`,
	`StringTraits\.h`, `string_traits.h`,
	`AutoMem\.h`, `auto_block.h`,
	`ReadIni\.h`, `ini_reader.h`,
	`IniFile\.h`, `ini_file.h`,
	`XmlFile\.h`, `xml_reader.h`,
}
var regfilerename RegRep

func init() {
	regkn = makere(kn)
	regku = makere(ku)
	reggu = makere(gu)
	regrepath = makere(repath)
	regfilerename = makere(filerename)
	regcbrename = makere(cbrename)
	regclassrename = makere(classRename)
	regrename = makere(rename)
	regvarlistupgrade = makere(varlistupgrade)
	regdefine = makere(redefine)
}

var reginclude = regexp.MustCompile(`#include`)
var regCamel = regexp.MustCompile(`"(?P<path>[\w\./]+)"`)
var regCamel2 = regexp.MustCompile(`\<(?P<path>[\w\./]+)\>`)
var regseparator = regexp.MustCompile(`\\`)

func GetParams(regEx *regexp.Regexp, url string) (paramsMap map[string]string) {
	match := regEx.FindStringSubmatch(url)

	paramsMap = make(map[string]string)
	for i, name := range regEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return
}

func includeUpgrade(line []byte) ([]byte, bool) {
	if !reginclude.Match(line) {
		return line, false
	}
	//分隔符修改
	line = regseparator.ReplaceAll(line, []byte{'/'})
	line = regrepath.Replace(line)
	line = regfilerename.Replace(line)

	path := GetParams(regCamel, string(line))
	if p, has := path["path"]; has {
		p = ReplaceUpper(p)
		line = regCamel.ReplaceAll(line, []byte(`"`+SnakeString(p)+`"`))
	}

	path2 := GetParams(regCamel2, string(line))
	if p, has := path2["path"]; has {
		p = ReplaceUpper(p)
		line = regCamel2.ReplaceAll(line, []byte(`<`+SnakeString(p)+`>`))
	}
	//转换成小写
	line = []byte(strings.ToLower(string(line)))
	return line, true
}

func kernelUpgrader(line []byte) []byte {
	line = regku.Replace(line)
	line = regkn.Replace(line)
	return line
}

func gameobjUpgrade(line []byte) []byte {
	line = reggu.Replace(line)
	return line
}

func varlistUpgrade(line []byte) []byte {
	line = regvarlistupgrade.Replace(line)
	return line
}

func renameUpgrade(line []byte) []byte {
	line = regrename.Replace(line)
	return line
}

func cbUpgrade(line []byte) []byte {
	line = regcbrename.Replace(line)
	return line
}

func classUpgrade(line []byte) []byte {
	line = regclassrename.Replace(line)
	return line
}

func defineUpgrade(line []byte) []byte {
	line = regdefine.Replace(line)
	return line
}

func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	pre := byte(0)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j && pre != '/' && pre != '\\' {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		pre = d
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func ReplaceUpper(base string) string {
	base = strings.Replace(base, "AIT", "ait", -1)
	base = strings.Replace(base, "AIR", "air", -1)
	base = strings.Replace(base, "AI", "ai", -1)
	base = strings.Replace(base, "VIP", "vip", -1)
	base = strings.Replace(base, "GMCC", "gmcc", -1)
	base = strings.Replace(base, "GM", "gm", -1)
	base = strings.Replace(base, "NPC", "npc", -1)
	return base
}

func UpgradeFile(f string) {
	file, err := os.Open(f)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	file.Close()
	//删除旧文件
	os.Remove(f)
	dir := filepath.Dir(f)
	base := getBaseName(f)
	base = ReplaceUpper(base)
	base = SnakeString(base)
	ext := filepath.Ext(f)
	newpath := dir + "/" + base + ext
	fmt.Println(newpath)
	if len(data) > 0 {
		Upgrade(data, newpath)
	}
}

var regprogram = regexp.MustCompile(`#program once`)

func getBaseName(f string) string {
	return strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
}

func Upgrade(filedata []byte, out string) {
	size := len(filedata)
	if filedata[size-1] != '\n' {
		filedata = append(filedata, '\n')
	}

	buf := bytes.NewBuffer(filedata)
	isheader := false
	outfile, _ := os.Create(out)
	writer := bufio.NewWriter(outfile)
	changeprogram := false
	def := fmt.Sprintf("_%s_%s_", strings.ToUpper(getBaseName(out)), strings.ToUpper(strings.Replace(filepath.Ext(out), ".", "", -1)))
	pre := size / 80 // 进度条刻度
	if pre == 0 {
		pre = 1
	}
	last := 0
	for {

		line, err := buf.ReadBytes('\n')
		if err != nil || io.EOF == err {
			break
		}
		if !changeprogram && regprogram.Match(line) {
			changeprogram = true
			line = []byte(fmt.Sprintf("#ifndef %s \n#define %s \n", def, def))
		} else {
			if line, isheader = includeUpgrade(line); !isheader {
				line = kernelUpgrader(line)
				line = gameobjUpgrade(line)
				line = varlistUpgrade(line)
				line = defineUpgrade(line)
				line = renameUpgrade(line)
				line = cbUpgrade(line)
				line = classUpgrade(line)
			}
		}
		writer.Write(line)
		p := (int)((buf.Len() - last) / pre)
		for i := 0; i < p; i++ {
			fmt.Print(".")
		}

		last = last + pre*p
	}
	if changeprogram {
		writer.WriteString(fmt.Sprintf("\n#endif // end of %s\n", def))
	}
	fmt.Println("[done]")
	writer.Flush()
	outfile.Close()
}
