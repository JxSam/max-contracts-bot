package notifier

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	executor "github.com/JxSam/max-contracts-bot/internal/bot/message"
	"github.com/JxSam/max-contracts-bot/internal/model"
	"github.com/JxSam/max-contracts-bot/internal/parser"
	"github.com/JxSam/max-contracts-bot/internal/repository"
)

type Contract struct {
	Id          int64
	CreatedAt   time.Time
	Description string
}

func (n *Notifier) updateContracts(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			contracts, err := n.contractsService.GetAllContractsService()
			if err != nil {
				n.log.Error("failed to get contracts", slog.Any("err", err))
				continue
			}

			for _, c := range contracts {
				contract := c

				event, err := parser.ParseLink(
					contract.Link,
					contract.Default,
				)
				fmt.Println(strconv.FormatInt(contract.Id, 10))
				if err != nil {
					n.log.Error(
						"failed to parse contract",
						slog.Any("contract_id", contract.Id),
						slog.Any("err", err),
					)
					continue
				}

				if !event.EventTime.After(contract.CreatedAt) {
					continue
				}

				n.log.Info(
					"contract updated",
					slog.Any("contract_id", contract.Id),
					slog.Any("old_date", contract.CreatedAt),
					slog.Any("new_date", event.EventTime),
				)

				err = n.contractsService.UpdateContractService(
					ctx,
					contract.Id,
					event.EventTime,
					event.Description,
					event.ItemName,
				)
				if err != nil {
					n.log.Error(
						"failed to update contract",
						slog.Any("contract_id", contract.Id),
						slog.Any("err", err),
					)
					continue
				}
				contract.CreatedAt = event.EventTime
				contract.Description = event.Description

				n.notifyContract(ctx, &contract, event, contract.LinkContract)
			}
		}
	}
}

func (n *Notifier) updateListContracts(ctx context.Context) {
	sources := map[string]string{
		"contract": "https://zakupki.gov.ru/epz/contract/search/rss?morphology=on&search-filter=Дате+размещения&fz44=on&contractStageList_0=on&contractStageList=0&budgetLevelsIdNameHidden=%7B%7D&customerIdOrg=03692000368%3AОБЛАСТНОЕ+ГОСУДАРСТВЕННОЕ+БЮДЖЕТНОЕ+УЧРЕЖДЕНИЕ+%22ЧЕЛЯБИНСКИЙ+РЕГИОНАЛЬНЫЙ+ЦЕНТР+НАВИГАЦИОННО-ИНФОРМАЦИОННЫХ+ТЕХНОЛОГИЙ%22zZ03692000368zZ03692000368zZzZ7453245467zZ-1zZ745301001zZ1127453008264&sortBy=UPDATE_DATE&pageNumber=1&sortDirection=false&recordsPerPage=_500&showLotsInfoHidden=false",
		"order":    "https://zakupki.gov.ru/epz/order/extendedsearch/rss.html?searchString=&morphology=on&savedSearchSettingsIdHidden=&exclTextHidden=&sortBy=UPDATE_DATE&pageNumber=1&sortDirection=false&recordsPerPage=_10&showLotsInfoHidden=false&fz44=on&fz223=on&af=on&ca=on&pc=on&pa=on&placingWayList=&selectedLaws=&etp=&nationalRegime=&bankSupportTag=&priceContractAdvantages44IdHidden=&priceContractAdvantages44IdNameHidden=%7B%7D&requirementsToPurchase44IdMap=&restrictionsToPurchase44=&priceFromGeneral=&priceFromGWS=&priceFromUnitGWS=&priceToGeneral=&priceToGWS=&priceToUnitGWS=&currencyIdGeneral=-1&currencyCodeGeneral=&advancePercentFrom=&advancePercentTo=&publishDateFrom=&publishDateTo=&updateDateFrom=&updateDateTo=&applSubmissionCloseDateFrom=&applSubmissionCloseDateTo=&EADateFrom=&EADateTo=&kbkChapter=&kbkSection=&kbkArticle=&kbkExpense=&customerIdOrg=03692000368%3AОБЛАСТНОЕ+ГОСУДАРСТВЕННОЕ+БЮДЖЕТНОЕ+УЧРЕЖДЕНИЕ+%22ЧЕЛЯБИНСКИЙ+РЕГИОНАЛЬНЫЙ+ЦЕНТР+НАВИГАЦИОННО-ИНФОРМАЦИОННЫХ+ТЕХНОЛОГИЙ%22zZ03692000368zZ03692000368zZzZ7453245467zZ-1zZ745301001zZ1127453008264&agencyIdOrg=&customerFz94id=&customerTitle=&specializedOrgIdOrg=&customerFz94id=&customerTitle=&customerPlace=&customerPlaceCodes=&oktmoIds=&oktmoIdsCodes=&delKladrIds=&delKladrIdsCodes=&deliveryPlaceGAR=&taxpayerCode=&okpd2Ids=&okpd2IdsCodes=&okpdIds=&okpdIdsCodes=&okdpIds=&okdpIdsCodes=&ktruCodeNameList=&ktruSelectedChcs=&ktruSelectedChcsNames=&ktruSelectedCharItemVersionIdList=&ktruSelectedRubricatorIdList=&ktruSelectedRubricatorName=&clItemsHiddenId=&clGroupHiddenId=&ktruSelectedPageNum=&okved2=&okved2Codes=&okvedIds=&okvedIdsCodes=&worktypeIds=&worktypeNames=&worktypeIdsParent=&selectedSubjectsIdHidden=&selectedSubjectsIdNameHidden=%7B%7D&kvrIds=&kvrIdsCodes=&koksIdsIdHidden=&koksIdsIdNameHidden=%7B%7D&participantName=&orderIKZInputNameYear=&orderIKZInputNameIkz=&orderIKZInputNameNumPZ=&orderIKZInputNameNumPGZ=&orderIKZInputNameOkpd2=&orderIKZInputNameKvr=&orderIKZInputNameRef=&mnnFarmNameIdMap=&tradeFarmNameIdMap=&medFormIdMap=&farmDosageIdMap=&gws=Выберите+тип+закупки&searchTextInAttachedFile=",
	}
	link_contract := ""
	ticker := time.NewTicker(50 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			baseContracts, err := n.contractsService.GetAllContractsService()
			if err != nil {
				n.log.Error("failed to get contracts", slog.Any("err", err))
				continue
			}

			baseMap := make(map[int64]repository.Contracts, len(baseContracts))
			for _, c := range baseContracts {
				baseMap[c.Id] = c
			}

			for key, source := range sources {
				var link string
				if key == "contract" {
					link = "https://zakupki.gov.ru/epz/contract/contractCard/rss?reestrNumber="
					link_contract = "https://zakupki.gov.ru/epz/contract/contractCard/common-info.html?reestrNumber="
				} else if key == "order" {
					link = "https://zakupki.gov.ru/epz/order/notice/rss?regNumber="
					link_contract = "https://zakupki.gov.ru/epz/order/notice/ea20/view/common-info.html?regNumber=0"
				}
				contracts, err := parser.ParseAllContracts(source)
				if err != nil {
					n.log.Error("failed to parse rss",
						slog.String("source", source),
						slog.Any("err", err),
					)
					continue
				}

				for _, contract := range contracts {
					contractID, err := strconv.ParseInt(contract, 10, 64)
					if err != nil {
						continue
					}

					// Уже есть в БД
					if baseContract, ok := baseMap[contractID]; ok {
						if !baseContract.Default {
							if err := n.contractsService.UpdateDefaultContract(ctx, contractID, true); err != nil {
								n.log.Error("failed to update default",
									slog.Int64("contract_id", contractID),
									slog.Any("err", err),
								)
							}
						}
						continue
					}

					event, err := parser.ParseLink(link+contract, false)
					if err != nil {
						n.log.Error("failed to parse contract",
							slog.Int64("contract_id", contractID),
							slog.Any("err", err),
						)
						continue
					}

					nowUTC := time.Now().UTC()
					eventTimeUTC := event.EventTime.UTC()
					defaultDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
					yesterdayStart := time.Now().UTC().AddDate(0, 0, -2)

					if !event.EventTime.After(yesterdayStart) {
						continue
					}

					if !eventTimeUTC.IsZero() && eventTimeUTC.After(nowUTC) {
						// Если дата в будущем — ставим дефолтную
						event.EventTime = defaultDate
						n.log.Debug("fixed future date",
							slog.Int64("contract_id", contractID),
							slog.Time("original_time", eventTimeUTC),
							slog.Time("fixed_time", defaultDate),
						)
					}

					err = n.contractsService.CreateContract(
						contractID,
						link+contract,
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						event.Description,
						event.ItemName,
						true,
						link_contract,
					)
					if err != nil {
						n.log.Error("failed to create contract",
							slog.Int64("contract_id", contractID),
							slog.Any("err", err),
						)
						continue
					}

					baseMap[contractID] = repository.Contracts{
						Id:      contractID,
						Default: true,
					}
				}
			}
		}
	}
}

func (n *Notifier) notifyContract(
	ctx context.Context,
	contract *repository.Contracts,
	event parser.Event,
	link string,
) {
	message := executor.Contract(contract)

	var userList []model.User

	if event.Default {
		users, err := n.userService.GetUsers()
		if err != nil {
			n.log.Error("failed to get users", slog.Any("err", err))
			return
		}

		userList = make([]model.User, 0, len(users))
		for _, u := range users {
			userList = append(userList, model.User{
				ChatID: u.ChatID,
			})
		}
	} else {
		users, err := n.active_users.GetUsersActiveContractsService(contract.Id)
		if err != nil {
			n.log.Error(
				"failed to get users_contracts",
				slog.Any("err", err),
			)
			return
		}

		userList = make([]model.User, 0, len(users))
		for _, u := range users {
			userList = append(userList, model.User{
				ChatID: u.ChatID,
			})
		}
	}

	msgBase := model.Notify{
		Message: *message,
		LinkButton: &model.LinkButton{
			Text: "Перейти к контракту",
			Link: link +
				strconv.FormatInt(contract.Id, 10),
		},
	}

	for _, u := range userList {
		msg := msgBase
		msg.UserID = u.ChatID

		select {
		case n.ch <- msg:
		case <-ctx.Done():
			return
		default:
			n.log.Warn(
				"notification channel full, skipping message",
				slog.Any("user_id", u.ChatID),
			)
		}
	}
}
