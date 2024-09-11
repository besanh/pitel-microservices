package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/tel4vn/pitel-microservices/common/cache"
	"github.com/tel4vn/pitel-microservices/common/log"
	"github.com/tel4vn/pitel-microservices/internal/sqlclient"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/tel4vn/pitel-microservices/repository"
	"gopkg.in/gomail.v2"
)

func (s *ChatEmail) HandleJobExpireToken() {
	dbCon, err := NewDBConn("", DBConfig{
		Host:     DB_HOST,
		Port:     DB_PORT,
		Database: DB_DATABASE,
		Username: DB_USERNAME,
		Password: DB_PASSWORD,
	})
	if err != nil {
		log.Error(err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := model.ChatConnectionAppFilter{
		Status: "active",
	}
	total, connections, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, dbCon, filter, -1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if total > 0 {
		for _, connection := range *connections {
			log.Info("start job processing connection: " + connection.ConnectionName)
			if len(connection.OaInfo.Facebook) > 0 {
				if connection.OaInfo.Facebook[0].Expire != 0 {
					if err := handleFlowExpireFacebook(ctx, dbCon, connection); err != nil {
						continue
					}
				}
			} else if len(connection.OaInfo.Zalo) > 0 {
				if connection.OaInfo.Zalo[0].Expire != 0 {
					if err := handleFlowExpireZalo(ctx, dbCon, connection); err != nil {
						continue
					}
				}
			}
		}
	}
}

/**
* Facebook expire in in 90 days ~ 7776000s
 */
func handleFlowExpireFacebook(ctx context.Context, dbCon sqlclient.ISqlClientConn, connection model.ChatConnectionApp) (err error) {
	// We can use updated_timestamp to compare
	expire := connection.OaInfo.Facebook[0].Expire
	if expire == 0 {
		createdTimestamp := connection.OaInfo.Facebook[0].CreatedTimestamp
		expire = createdTimestamp + 7776000
	}
	if expire == 7776000 {
		return errors.New("expire time " + fmt.Sprintf("%d", expire) + " is not valid")
	}
	expireTime, err := ParseExpire(expire)
	if err != nil {
		log.Error(err)
		return
	}

	/*
	* TODO: compare (expire time token - 1 day) with current time
	* From x-(1day + 1hour) to x
	 */
	dayBeforeStartExpire := expireTime.AddDate(0, 0, -1)
	dayBeforeEndExpire := expireTime.Add(-90000)

	if time.Now().After(dayBeforeStartExpire) && time.Now().Before(dayBeforeEndExpire) && !connection.OaInfo.Facebook[0].IsNotify {
		// Caching connection
		chatEmail := model.ChatEmail{}
		connectionCache := cache.RCache.Get(CHAT_CONNECTION + "_" + connection.OaInfo.Facebook[0].AppId + "_" + "_" + connection.OaInfo.Facebook[0].OaId)
		if connectionCache != nil {
			if err = json.Unmarshal(connectionCache.([]byte), &chatEmail); err != nil {
				log.Error(err)
				return
			}
		} else {
			filter := model.ChatEmailFilter{
				OaId: connection.OaInfo.Facebook[0].OaId,
			}
			_, emails, err := repository.NewChatEmail().GetChatEmails(ctx, dbCon, filter, 1, 0)
			if err != nil {
				log.Error(err)
				return err
			} else if len(*emails) < 1 {
				log.Error("email " + connection.OaInfo.Facebook[0].OaId + " not found")
				return err
			}
			chatEmail = (*emails)[0]

			// TODO: set cache
			// if err = cache.RCache.Set(CHAT_CONNECTION+"_"+connection.OaInfo.Facebook[0].OaId, chatEmail, 5*time.Minute); err != nil {
			// 	log.Error(err)
			// 	return err
			// }
		}

		// TODO: send email
		if err = SendEmailAsync(chatEmail); err != nil {
			log.Error(err)

			if SMTP_INFORM {
				// Send mail to noreply to notify error
				if errInform := SendEmailAsync(model.ChatEmail{
					EmailUsername: SMTP_USERNAME,
					EmailRecipient: []string{
						SMTP_USERNAME,
					},
					EmailSubject: "Error sending email - Oa type: " + connection.ConnectionType + " - Oa id: " + connection.OaInfo.Facebook[0].OaId + " - Oa name: " + connection.OaInfo.Facebook[0].OaName,
					EmailContent: err.Error(),
				}); errInform != nil {
					log.Error(errInform)
					return errInform
				}
			}

			return err
		}

		// Update connection mark notified
		connection.OaInfo.Facebook[0].IsNotify = true
		if err = repository.ChatConnectionAppRepo.Update(ctx, dbCon, connection); err != nil {
			log.Error(err)
			return err
		}
	}

	return
}

/**
* Zalo expire in in 30 days
 */
func handleFlowExpireZalo(ctx context.Context, dbCon sqlclient.ISqlClientConn, connection model.ChatConnectionApp) (err error) {
	if connection.OaInfo.Zalo[0].UpdatedTimestamp == 0 {
		return errors.New("expire time " + fmt.Sprintf("%d", connection.OaInfo.Zalo[0].UpdatedTimestamp) + " is not valid")
	}
	expireTime, err := ParseExpire(connection.OaInfo.Zalo[0].UpdatedTimestamp + connection.OaInfo.Zalo[0].TokenTimeRemainning)
	if err != nil {
		log.Error(err)
		return
	}

	/*
	* TODO: compare (expire time token - 1 day) with current time
	* From x-(29days) to x
	 */
	dayBeforeStartExpire := expireTime.AddDate(0, 0, -1)
	dayBeforeEndExpire := expireTime.Add(-90000)

	if time.Now().After(dayBeforeStartExpire) && time.Now().Before(dayBeforeEndExpire) && !connection.OaInfo.Zalo[0].IsNotify {
		// Caching connection
		chatEmail := model.ChatEmail{}
		connectionCache := cache.RCache.Get(CHAT_CONNECTION + "_" + connection.TenantId + "_" + connection.OaInfo.Zalo[0].AppId + "_" + connection.OaInfo.Zalo[0].OaId)
		if connectionCache != nil {
			if err = json.Unmarshal(connectionCache.([]byte), &chatEmail); err != nil {
				log.Error(err)
				return
			}
		} else {
			filter := model.ChatEmailFilter{
				OaId: connection.OaInfo.Zalo[0].OaId,
			}
			_, emails, err := repository.NewChatEmail().GetChatEmails(ctx, dbCon, filter, 1, 0)
			if err != nil {
				log.Error(err)
				return err
			} else if len(*emails) < 1 {
				log.Error("email " + connection.OaInfo.Zalo[0].OaId + " not found")
				return err
			}
			chatEmail = (*emails)[0]

			// TODO: set cache
			// if err = cache.RCache.Set(CHAT_CONNECTION+"_"+connection.OaInfo.Zalo[0].OaId, chatEmail, 5*time.Minute); err != nil {
			// 	log.Error(err)
			// 	return err
			// }
		}

		// TODO: send email
		if err = SendEmailAsync(chatEmail); err != nil {
			log.Error(err)

			if SMTP_INFORM {
				// Send mail to noreply to notify error
				if errInform := SendEmailAsync(model.ChatEmail{
					EmailUsername: SMTP_USERNAME,
					EmailRecipient: []string{
						SMTP_USERNAME,
					},
					EmailSubject: "Error sending email - Oa type: " + connection.ConnectionType + " - Oa id: " + connection.OaInfo.Zalo[0].OaId + " - Oa name: " + connection.OaInfo.Zalo[0].OaName,
					EmailContent: err.Error(),
				}); errInform != nil {
					log.Error(errInform)
					return errInform
				}
			}

			return err
		}

		// Update connection mark notified
		connection.OaInfo.Zalo[0].IsNotify = true
		if err = repository.ChatConnectionAppRepo.Update(ctx, dbCon, connection); err != nil {
			log.Error(err)
			return err
		}
	}

	return
}

func ParseExpire(expire int64) (expireTime time.Time, err error) {
	i, err := strconv.ParseInt(fmt.Sprintf("%d", expire), 10, 64)
	if err != nil {
		return expireTime, err
	}
	expireTime = time.Unix(i, 0)
	return
}

func SendEmailAsync(email model.ChatEmail) (err error) {
	mail := gomail.NewMessage()
	mail.SetHeaders(map[string][]string{
		"From":    {email.EmailUsername},
		"To":      email.EmailRecipient,
		"Subject": {email.EmailSubject},
	})

	mail.SetBody("text/html", email.EmailContent)

	emailPort, _ := strconv.Atoi(email.EmailPort)
	dialer := gomail.NewDialer(email.EmailServer, emailPort, email.EmailUsername, email.EmailPassword)
	sender, err := dialer.Dial()
	if err != nil {
		log.Error(err)
		return
	}
	defer sender.Close()
	if err = sender.Send(email.EmailUsername, email.EmailRecipient, mail); err != nil {
		log.Error(err)
		return
	}

	return
}
