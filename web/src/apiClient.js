// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

let cachedInfo = {};
let cachedProperties = {};

export function Post(url, body) {
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

export async function Login(username, password) {
	const res = await Post('/api/v0/authenticate', {
		Username: username,
		Password: password
	});

	if (res.ok) {
		cachedInfo = {};
		cachedProperties = {};
	}
	
	return Promise.resolve(res.ok);
}

export function Logout() {
	cachedInfo = {};
	cachedProperties = {};

	return fetch('/api/v0/logout');
}

export async function GetInfo() {
	if (cachedInfo.Authenticated !== undefined) {
		return cachedInfo;
	}

	else {
		cachedInfo = await fetch('/api/v0/info').then(res => res.json());
		return cachedInfo;
	}
}

export async function GetPool(id) {
	return await fetch('/api/v0/pool?pool=' + id).then(res => res.json());
}

export async function GetFields(table) {
	if (cachedProperties[table] !== undefined) {
		return cachedProperties[table];
	}

	return await fetch('/api/v0/properties?type=' + table)
	.then(res => {
		cachedProperties[table] = res.json();
		return cachedProperties[table];
	});
}

async function getText(url) {
	return await fetch(url)
	.then((res) => {
		return Promise.resolve(res.text());
	});
}

export async function Mount(id) {
	return await getText('/api/v0/data/mount?id=' + id);
}

export async function Unmount(id) {
	return await getText('/api/v0/data/unmount?id=' + id);
}

export async function LoadKey(id, passphrase) {
	const res = await Post('/api/v0/key/load', {
		'id': id,
		'passphrase': passphrase
	});

	return await res.text();
}

export async function UnloadKey(id) {
	const res = await Post('/api/v0/key/unload', {
		'id': id
	});

	return await res.text();
}