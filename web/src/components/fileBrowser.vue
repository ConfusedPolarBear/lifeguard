<template><div>
    <b-modal id="modalBrowser" centered size="xl">
        <p>Current path: <code>{{ path }}</code></p>

        <b-table :items="contents" :fields="fields">
            <template v-slot:cell(name)="data">
                <span class="material-icons">{{ typeToIcon(data.item.Type) }}</span>

                <span v-if="data.item.Type == 'd'">
                    <b-link href="#" :id="data.item.HMAC" :data-type="data.item.Type" @click="loadEntry">{{ data.item.Name }}</b-link>
                </span>
                <span v-else>
                    <b-link :href="'/api/v0/files/browse/' + data.item.HMAC" target="_blank">{{ data.item.Name }}</b-link>
                </span>
            </template>

            <template v-slot:cell(size)="data">
                {{ data.item.Size | prettyPrint('size') }}
            </template>
        </b-table>
    </b-modal>
</div></template>

<script>
import * as ApiClient from '../apiClient.js';

export default {
	name: 'file-browser',
	props: [ 'name', 'hmac' ],
	data() {
		return {
			'path': '',
			'contents': {},
			'fields': [
				{
					key: 'Name',
					sortable: true
				},
				{
					key: 'Size',
					sortable: true
				}
			]
		};
	},
	mounted: function() {
		this.browse(this.hmac);
		this.$bvModal.show('modalBrowser');
	},
	methods: {
		browse: async function(hmac) {
			this.contents = await ApiClient.Browse(hmac);
			this.path = this.contents[0].Name;

			// pop the zero'th element off - this is the current path
			this.contents.shift();
		},
		loadEntry: async function(e) {
			console.debug(e);

			let hmac = e.target.id;
			let type = e.target.dataset.type;

			if (type === 'd') {
				this.browse(hmac);
			}
		},
		typeToIcon: function(type) {
			return (type === 'd') ? 'folder' : 'description';
		}
	}
};
</script>