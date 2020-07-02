// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

let cachedInfo = {};

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

async function Login(username, password) {
	const res = await Post('/api/v0/authenticate', {
		Username: username,
		Password: password
	});

	if (res.ok) {
		cachedInfo = {};
	}
	
	return Promise.resolve(res.ok);
}

function Logout() {
	cachedInfo = {};
	return fetch('/api/v0/logout');
}

async function GetInfo() {
	if (cachedInfo.Authenticated !== undefined) {
		return cachedInfo;
	}

	else {
		cachedInfo = await fetch('/api/v0/info').then(res => res.json());
		return cachedInfo;
	}
}

async function GetFields(table) {
	return await fetch('/api/v0/properties?type=' + table)
		.then(res => res.json());
}

export default {
	Post: Post,
	Login: Login,
	Logout: Logout,
	GetInfo: GetInfo,
	GetFields: GetFields
}
