// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

import * as ApiClient from '../apiClient.js';

export async function Trim(pool) {
	const res = await ApiClient.Post('/api/v0/pool/' + pool + '/trim', {});
	return await res.text();
}

export async function Iostat(pool) {
	const res = await fetch('/api/v0/pool/' + pool + '/iostat', {});
	if (!res.ok) {
		return Promise.reject(await res.text());
	}

	return Promise.resolve(await res.text());
}