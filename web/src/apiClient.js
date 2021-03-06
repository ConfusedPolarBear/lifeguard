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

	if (!res.ok) {
		return Promise.reject('Bad username/password');
	}

	let authType = await res.text();
	authType = authType.replaceAll('\n', '');

	if (authType === 'full') {
		cachedInfo = {};
		cachedProperties = {};
		console.log('login successful');

	} else if (authType === 'partial') {
		console.log('partially logged in, waiting for response to 2FA challenge');

	} else {
		console.log('unknown auth result "' + authType + '"');
		Promise.reject('Unknown response ' + authType);
	}
	
	return Promise.resolve(authType);
}

export function Logout() {
	cachedInfo = {};
	cachedProperties = {};

	return Post('/api/v0/logout');
}

export async function GetInfo() {
	if (cachedInfo.Authenticated !== undefined) {
		return cachedInfo;
	} else {
		let info = await fetch('/api/v0/info').then(res => res.json());
		
		// Never cache info from before we are authenticated
		if (info.Authenticated) {
			console.log('caching info');
			cachedInfo = info;
		} else {
			console.log('not caching info');
		}

		return info;
	}
}

export async function GetPool(id) {
	const res = await fetch('/api/v0/pool/' + encodeURIComponent(id));
	return await res.json();
}

export async function GetPools() {
	const res = await fetch('/api/v0/pools');
	return await res.json();
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

export async function Browse(id) {
	const res = await fetch('/api/v0/files/browse/' + encodeURIComponent(id));
	return await res.json();
}

export async function Scrub(id) {
	const res = await Post('/api/v0/pool/' + encodeURIComponent(id) + '/scrub/start');
	return await res.text();
}

export async function PauseScrub(id) {
	const res = await Post('/api/v0/pool/' + encodeURIComponent(id) + '/scrub/pause');
	return await res.text();
}

export async function GetNotifications() {
	const res = await fetch('/api/v0/notifications/list');
	
	let list = await res.json();
	return list.reverse();
}

export async function GetTwoFactorChallenge() {
	const res = await fetch('/api/v0/tfa/challenge');
	return await res.json();
}