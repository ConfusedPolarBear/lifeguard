// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

import * as ApiClient from '../apiClient.js';

export async function Initialize() {
	const res = await fetch('/api/v0/tfa/totp/initialize');
	return await res.json();
}

export async function IsEnabled() {
	const res = await fetch('/api/v0/tfa/enabled');
	let json = await res.json();
	return json.Enabled;
}

export async function Save(secret, code) {
	const res = await ApiClient.Post('/api/v0/tfa/totp/save', {
		secret: secret,
		code: code
	});

	if (!res.ok) {
		return Promise.reject(await res.text());
	}

	return Promise.resolve(true);
}

export async function Authenticate(code) {
	const res = await ApiClient.Post('/api/v0/tfa/totp/authenticate', {
		code: code
	});

	return Promise.resolve(res.ok);
}
