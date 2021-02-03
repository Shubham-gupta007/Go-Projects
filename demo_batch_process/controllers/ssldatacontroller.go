package controllers

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"zona_tip_batch/tip_batch_lib/common"
	"zona_tip_batch/tip_batch_lib/config"
	securaaelastic "zona_tip_batch/tip_batch_lib/elastic"
	itype "zona_tip_batch/tip_batch_lib/itypeconstants"
	"zona_tip_batch/tip_batch_lib/logger"
	"zona_tip_batch/tip_batch_lib/models"
	"zona_tip_batch/tip_services/utils"

	"zona_tip_batch/tip_batch_abuse.ch/constants"

	elastic "github.com/olivere/elastic"
)

type (
	SSLDataController struct{}
)

func NewSSLDataController() *SSLDataController {
	return &SSLDataController{}
}

func (ssl SSLDataController) GetSSLInformation(configObject config.ConfigStruct, client *elastic.Client) {
	logObj := logger.ErrorLog{}
	logObj.Source = "abusech"
	logObj.LogFileName = "ssldatacontroller.go"
	logObj.LogFunctionName = "GetSSLInformation"

	data, responseCode := common.GetRestData(constants.Abuse_sslblocklist, logObj)
	if responseCode == 200 {
		resp := strings.Split(data, "\n")
		r := csv.NewReader(strings.NewReader(data))

		bulkRequest := client.Bulk()
		req := true
		count := 0
		indicatorcount := 0
		var indicatorsMap []string
		syncID := "sid_" + strconv.FormatInt(common.GetCurrentTimestamp(), 10)

		logObj.BatchLog("Data Count from Response")
		logObj.BatchLog(common.ToJsonString(len(resp)))
		logObj.BatchLog("Inserting/Updating Data to ElasticSearch..")

		// for _, eachRecord := range resp {
		// 	if strings.HasPrefix(eachRecord, "#") {
		// 		continue
		// 	}
		for {
			eachRecord, err := r.Read()
			if err == io.EOF {
				break
			}

			// eachRecord = strings.ReplaceAll(eachRecord, "\r", "")
			// columnData := strings.Split(eachRecord, ",")
			if len(eachRecord) > 1 {
				if !strings.HasPrefix(eachRecord[0], "#") {
					indicator := eachRecord[1]

					if len(indicator) == 0 {
						logObj.BatchLog("Invalid Indicator")
						continue
					}

					exists, _ := utils.In_array(indicator, indicatorsMap)
					if exists {
						logObj.BatchLog("Indicator " + indicator + " already exist in current queue")
						continue
					}
					indicatorcount++
					count++

					indicatorsMap = append(indicatorsMap, indicator)
					ttlTimestamp := common.GetTTLTimestamp()
					processedTS := common.GetCurrentTimestamp()
					firstSeen := common.GetUnixTimestamp(eachRecord[0])

					eachColumn := models.IndicatorData{
						Indicator:   indicator,
						Itype:       itype.SSL,
						Ttl:         ttlTimestamp,
						Source:      constants.SourceSSL,
						SourceLink:  constants.Abuse_sslblocklist,
						Reason:      eachRecord[2],
						GeoFlag:     constants.DefaultFlag,
						AsnFlag:     constants.DefaultFlag,
						FirstSeen:   firstSeen,
						Confidence:  constants.Confidence,
						UpdatedTS:   processedTS,
						ProcessedTS: processedTS,
						SyncID:      syncID,
						UUID:        common.GetIndicatorUUID()}

					req, bulkRequest = securaaelastic.UpdateBulkRequestObject(configObject.ElasticSearchHost, configObject.ElasticSearchPort, eachColumn, indicator, constants.SourceSSL, bulkRequest, logObj)
					if !req {
						break
					} else {
						if (count % 250) == 0 {
							count = 0
							indicatorsMap = nil
							_, err := securaaelastic.ExecuteBulkRequest(bulkRequest, logObj)
							if err != nil {
								logObj.BatchLog("Error 1 : " + err.Error())
							} else {
								bulkRequest = nil
								bulkRequest = client.Bulk()
								logObj.BatchLog("Executed BulkRequest")
							}
						}
					}
				}
			}
		}

		if count > 0 {
			_, err := securaaelastic.ExecuteBulkRequest(bulkRequest, logObj)
			if err != nil {
				logObj.BatchLog("Error 2 : " + err.Error())
			} else {
				logObj.BatchLog("Executed BulkRequest")
			}
		}

		indicatorsMap = nil
		logObj.BatchLog("Inserted/Updated rows count " + common.ToJsonString(indicatorcount) + "/" + common.ToJsonString(len(resp)))

		// UPDATE SYNC DATA COUNT AND TIME
		securaaelastic.UpdateDataSync(constants.SourceSSL, indicatorcount, logObj, syncID)
	} else {
		logObj.BatchLog("-----------------------------")
		logObj.BatchLog("NO DATA FROM " + constants.Abuse_sslblocklist + ". RESPONSE CODE : " + strconv.Itoa(responseCode))
		logObj.BatchLog("-----------------------------")
	}
}

