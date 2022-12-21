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
	UserID       string `form:"user_id" json:"user_id" xml:"user_id"`
	Username     string `form:"username" json:"username" xml:"username" validation:"username"`
	ClientID     string `form:"client_id" json:"client_id" xml:"client_id" validation:"client_id"`
	RedirectUri  string `form:"redirect_uri" json:"redirect_uri" xml:"redirect_uri" validation:"redirect_uri"`
	ResponseType string `form:"response_type" json:"response_type" xml:"response_type" validation:"response_type"`
	Scopes       string `form:"scopes" json:"scopes" xml:"scopes" validation:"scopes"`
	PubKey       string `form:"pub_key" json:"pub_key" xml:"pub_key" validation:"pub_key"`
}
