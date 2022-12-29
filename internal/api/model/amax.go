package model

/*
   GoToSocial
   Copyright (C) 2021-2022 GoToSocial Authors admin@gotosocial.org

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

type AmaxSubmitInfoRequest struct {
	ClientName   string `form:"client_name" json:"client_name" xml:"client_name" validation:"client_name"`
	RedirectUris string `form:"redirect_uris" json:"redirect_uris" xml:"redirect_uris" validation:"redirect_uris"`
	Scope        string `form:"scope" json:"scope" xml:"scope" validation:"scope"`
	GrantType    string `form:"grant_type" json:"grant_type" xml:"grant_type" validation:"grant_type"`
	ClientId     string `form:"client_id" json:"client_id" xml:"client_id" validation:"client_id"`
	ClientSecret string `form:"client_secret" json:"client_secret" xml:"client_secret" validation:"client_secret"`
	Reason       string `form:"reason" json:"reason" xml:"reason" validation:"reason"`
	Email        string `form:"email" json:"email" xml:"email" validation:"email"`
	Username     string `form:"username" json:"username" xml:"username" validation:"username"`
	Password     string `form:"password" json:"password" xml:"password" validation:"password"`
	Agreement    bool   `form:"agreement" json:"agreement" xml:"agreement" validation:"agreement"`
	Locale       string `form:"locale" json:"locale" xml:"locale" validation:"locale"`
}

type AmaxSignatureLoginRequest struct {
	Username string `form:"username" json:"username" xml:"username" validation:"username"`
	PubKey   string `form:"pub_key" json:"pub_key" xml:"pub_key" validation:"pub_key"`
}

type AmaxSignatureRequest struct {
	Signature  string `form:"signature" json:"signature" xml:"signature" validation:"signature"`
	Address    string `form:"address" json:"address" xml:"address" validation:"address"`
	ChainId    string `form:"chainId" json:"chainId" xml:"chainId" validation:"chainId"`
	Authority  string `form:"authority" json:"authority" xml:"authority" validation:"authority"`
	Expiration string `form:"expiration" json:"expiration" xml:"expiration" validation:"expiration"`
	Scope      string `form:"scope" json:"scope" xml:"scope" validation:"scope"`
}
