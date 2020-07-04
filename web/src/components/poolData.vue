<!--
Copyright 2020 Matt Montgomery
SPDX-License-Identifier: AGPL-3.0-or-later
-->

<template><div class="poolData">
    <b-card :header="section">
        <b-form-group>
            <input style="width:45%;margin-right:1em" type="text" placeholder="Filter" v-model="filter">

            <b-button :disabled="disableToolbar" @click="mount">{{ propertyEqual('mounted', 'no') ? 'Mount': 'Unmount' }}</b-button>
            <b-button :disabled="disableToolbar" @click="loadKey">{{ propertyEqual('keystatus', 'unavailable') ? 'Load key' : 'Unload key' }}</b-button>
            <b-dropdown split text="Snapshot"  :disabled="disableToolbar">
                <b-dropdown-item>Diff</b-dropdown-item>
                <b-dropdown-divider></b-dropdown-divider>
                <b-dropdown-item variant="danger">Rollback</b-dropdown-item>
                <b-dropdown-item variant="danger">Prune</b-dropdown-item>
            </b-dropdown>
        </b-form-group>

        <b-table selectable select-mode="range" @row-selected="onSelect" outlined hover :fields="fields[section]" :items="pool[section]" :filter="filter">
            <template v-slot:cell()="data">
                {{ data.value.Value | prettyPrint(data.value.Name) }}
            </template>
        </b-table>
    </b-card>
</div></template>

<script>
export default {
    name: 'poolData',
    props: [ 'pool', 'section', 'fields' ],
    data() { return {
        filter: '',
        selected: [],
    } },
    computed: {
        disableToolbar: function() {
            return this.selected.length === 0;
        }
    },
    methods: {
        onSelect: function(items) {
            this.selected = items;

            // Bug fix: If a filter is active and a selection is made, the selection will be cleared on the next background refresh
            this.$emit('select', items, this.filter);
        },
        propertyEqual: function(prop, state) {
            if (this.selected.length > 1) {
                // TODO: implement this for multiple datasets
                return false;
            }

            else if (this.selected.length === 0) {
                return true;
            }

            return this.selected[0][prop].Value === state;
        },
        mount: function() {
            let event = this.propertyEqual('mounted', 'no') ? 'mount' : 'unmount';
            this.$emit('click', event, this.selected[0].name.Value);
        },
        loadKey: function() {
            let event = this.propertyEqual('keystatus', 'unavailable') ? 'load-key' : 'unload-key';
            this.$emit('click', event, this.selected[0].name.Value);
        }
    }
}
</script>