package user

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"git.friends.com/PetLand/UserService/v2/internal/core"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"git.friends.com/PetLand/UserService/v2/internal/models"
	"git.friends.com/PetLand/UserService/v2/internal/store"
	"gopkg.in/gomail.v2"
)

func ForgotPassword(data *models.Contacts, email string, store *store.Store) error {
	//code := "1234"
	hash := sha256.New()
	hash.Write([]byte(data.ProfileID.String()))
	hashSend := hex.EncodeToString(hash.Sum(nil))
	fmt.Printf("%s\n", hashSend)
	//a := string(hashSend)
	err := store.Contacts().InsertHashID(data.ProfileID, hashSend)
	if err != nil {
		return genErr.NewError(err, core.ErrRepository, "msg", "failed insert hash")
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", "petland23@mail.ru")
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Восстановление пароля")
	msg.SetBody("text/html", fmt.Sprintf("<!DOCTYPE html\n        PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n<html xmlns=\"http://www.w3.org/1999/xhtml\" xmlns:v=\"urn:schemas-microsoft-com:vml\"\n      xmlns:o=\"urn:schemas-microsoft-com:office:office\" lang=\"ru\">\n\n<head>\n    <meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\" />\n    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\" />\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\n    <meta name=\"color-scheme\" content=\"light dark\" />\n    <meta name=\"supported-color-schemes\" content=\"light dark\" />\n    <title>PetLandMail</title>\n    <style type=\"text/css\">\n        @import url('https://fonts.googleapis.com/css2?family=Modak&family=Mulish:wght@400;700&display=swap');\n\n        table {\n            border-spacing: 0;\n            mso-cellspacing: 0;\n            mso-padding-alt: 0;\n        }\n\n        td {\n            padding: 0;\n        }\n\n        #outlook a {\n            padding: 0;\n        }\n\n        a {\n            text-decoration: none;\n            color: #e8fbfa;\n            font-size: 16px;\n        }\n\n        @media screen and (max-width: 599.98px) {}\n\n        @media screen and (max-width: 399.98px) {\n            .mobile-padding {\n                padding-right: 10px !important;\n                padding-left: 10px !important;\n            }\n\n            .mobile-col-padding {\n                padding-right: 0 !important;\n                padding-left: 0 !important;\n            }\n\n            .two-columns .column {\n                width: 100%% !important;\n                max-width: 100%% !important;\n            }\n\n            .two-columns .column img {\n                width: 100%% !important;\n                max-width: 100%% !important;\n            }\n\n            .three-columns .column {\n                width: 100%% !important;\n                max-width: 100%% !important;\n            }\n\n            .three-columns .column img {\n                width: 100%% !important;\n                max-width: 100%% !important;\n            }\n        }\n\n        /* Custom Dark Mode Colors */\n        :root {\n            color-scheme: light dark;\n            supported-color-schemes: light dark;\n        }\n\n        @media (prefers-color-scheme: dark) {\n\n            table,\n            td {\n                background-color: #06080B !important;\n            }\n\n            h1,\n            h2,\n            h3,\n            p {\n                color: #ffffff !important;\n            }\n        }\n    </style>\n\n    <!--[if (gte mso 9)|(IE)]>\n    <style type=\"text/css\">\n        table {border-collapse: collapse !important;}\n    </style>\n    <![endif]-->\n\n    <!--[if (gte mso 9)|(IE)]>\n    <xml>\n        <o:OfficeDocumentSettings>\n            <o:AllowPNG/>\n            <o:PixelsPerInch>96</o:PixelsPerInch>\n        </o:OfficeDocumentSettings>\n    </xml>\n    <![endif]-->\n</head>\n\n<body style=\"Margin:0;padding:0;min-width:100%%;background-color:#dde0e1;\">\n\n<!--[if (gte mso 9)|(IE)]>\n<style type=\"text/css\">\n    body {background-color: #dde0e1!important;}\n    body, table, td, p, a {font-family: sans-serif, Arial, Helvetica!important;}\n</style>\n<![endif]-->\n\n<center style=\"width: 100%%;table-layout:fixed;background-color: #fff;padding-top: 40px;padding-bottom: 40px;\">\n    <div style=\"max-width: 600px;background-color: #fafdfe;box-shadow: 0 0 10px rgba(0, 0, 0, .2);\">\n\n        <!-- Preheader (remove comment) -->\n        <div\n                style=\"font-size: 0px;color: #fafdfe;line-height: 1px;mso-line-height-rule:exactly;display: none;max-width: 0px;max-height: 0px;opacity: 0;overflow: hidden;mso-hide:all;\">\n            Восстановление пароля.\n        </div>\n        <!-- End Preheader (remove comment) -->\n\n        <!--[if (gte mso 9)|(IE)]>\n        <table width=\"600\" align=\"center\" border=\"0\" cellspacing=\"0\" cellpadding=\"0\" role=\"presentation\"\n               style=\"color:#1C1E23;\">\n            <tr>\n                <td>\n        <![endif]-->\n\n        <table align=\"center\" border=\"0\" cellspacing=\"0\" cellpadding=\"0\" role=\"presentation\"\n               style=\"color:#4F4F4F;font-family: 'Mulish',sans-serif;background: #F5F1EE url('https://i.ibb.co/86RQrg3/bg-pattern.png');;Margin:0;padding:0;width: 100%%;max-width: 600px;\"\n        >\n\n\n            <!-- Hero -->\n            <tr>\n                <td style=\"padding: 0px 24px 25px 24px;\">\n                    <table border=\"0\" cellspacing=\"0\" cellpadding=\"0\" role=\"presentation\" style=\"width: 100%%; max-width: 600px;\">\n                        <tr>\n                            <td style=\"padding: 100px  0 0 0;\">\n                                <img src=\"https://i.ibb.co/z4YNQS0/PetLand.png\" alt=\"logo\" style=\"width: 377px; margin: 0 auto; display: block;\"/>\n                            </td>\n                        </tr>\n                        <tr>\n                            <td style=\"padding: 0 0 25px 0;\">\n                                <p style=\"margin: 0;font-weight: 400; font-size: 24px; text-align: center\">С заботой о ваших питомцах</p>\n                            </td>\n                        </tr>\n                        <tr>\n                            <td style=\"padding: 0 0 100px 0;\">\n                                <div style=\"margin: 0 auto;font-weight: 700;font-size: 16px;line-height: 25px;color: #1C1E23; width: 500px; background-color: white; border-radius: 20px; padding: 50px 0 50px 0; text-align: center\">Перейдите по ссылке для смены пароля - <span><a href=\"http://raezhov.fun/new-password?id=%s\" style=\"color:#98B14B; \">http://raezhov.fun/new-password?id=%s</a></span></div>\n                            </td>\n                        </tr>\n\n                    </table>\n                </td>\n            </tr>\n            <!-- End Hero -->\n\n        </table>\n        <![endif]-->\n\n    </div>\n</center>\n\n</body>\n\n</html>", hashSend, hashSend))

	n := gomail.NewDialer("smtp.mail.ru", 465, "petland23@mail.ru", "uCjsve57KRhjdik4FUBt")
	//n.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		return genErr.NewError(err, core.ErrWrongType)
	}

	return nil

}
