// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

function Post(url, body) {
	// Convert a raw object into a multipart form
	var data = new FormData();
	for (const key in body) {
		data.append(key, body[key]);
	}

	return fetch(url, {
		method: 'POST',
		body: data
	});
}

function Login(username, password) {
	return Post('/api/v0/authenticate', {
		Username: username,
		Password: password
	}).then(res => {
		return Promise.resolve(res.ok);
	});
}

function Logout() {
	return fetch('/api/v0/logout');
}

function GetInfo() {
	return fetch('/api/v0/info')
		.then(res => res.json());
}

export default {
	Post: Post,
	Login: Login,
	Logout: Logout,
	GetInfo: GetInfo
}
