// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

let cachedInfo = {};
let cachedProperties = {};

export function Post(url, body) {
	// Transform the object into a URI encoded form (the hard way)
	let data = [];
	for (let key in body) {
		data.push(encodeURIComponent(key) + '=' + encodeURIComponent(body[key]));
	}
	let postData = data.join('&');

	// Send it
	return fetch(url, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/x-www-form-urlencoded'
		},
		body: postData
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

	return Post('/api/v0/logout');
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
	return await fetch('/api/v0/pool/' + encodeURIComponent(id)).then(res => res.json());
}

export async function GetFields(table) {
	if (cachedProperties[table] !== undefined) {
		return cachedProperties[table];
	}

	return await fetch('/api/v0/properties/' + table)
	.then(res => {
		cachedProperties[table] = res.json();
		return cachedProperties[table];
	});
}

export async function GetSupportBundle() {
	return await getText('/api/v0/support');
}

async function getText(url) {
	return await fetch(url)
	.then((res) => {
		return Promise.resolve(res.text());
	});
}

// TODO: is there a better way to encode these URLs?
export async function Mount(id) {
	const res = await Post('/api/v0/data/' + encodeURIComponent(id) + '/mount');
	return await res.text();
}

export async function Unmount(id) {
	const res = await Post('/api/v0/data/' + encodeURIComponent(id) + '/unmount');
	return await res.text();
}

export async function LoadKey(id, passphrase) {
	const res = await Post('/api/v0/key/' + encodeURIComponent(id) + '/load', {
		'id': id,
		'Passphrase': passphrase
	});

	return await res.text();
}

export async function UnloadKey(id) {
	const res = await Post('/api/v0/key/' + encodeURIComponent(id) + '/unload', {
		'id': id
	});

	return await res.text();
}