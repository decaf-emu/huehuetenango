<template>
  <div class="uk-overflow-auto uk-padding-small uk-padding-remove-horizontal">
    <table class="uk-table uk-table-striped uk-table-small" v-if="functions && functions.length">
      <caption>Functions</caption>
      <thead>
        <tr>
          <th>Name</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in functions" :key="item">
          <td>
            <code>
              {{ item }}
            </code>
          </td>
        </tr>
      </tbody>
    </table>
    <table class="uk-table uk-table-striped uk-table-small" v-if="data && data.length">
      <caption>Data</caption>
      <thead>
        <tr>
          <th>Name</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in data" :key="item">
          <td>
            <code>
              {{ item }}
            </code>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  computed: {
    ...mapGetters(['exportData', 'exportFunctions']),
  	data() {
      if (!this.exportData) {
        return [];
      }

    	return this.exportData.map(item => Object.freeze(item));
    },
  	functions() {
      if (!this.exportFunctions) {
        return [];
      }

    	return this.exportFunctions.map(item => Object.freeze(item));
    },
  },
};
</script>

<style>
</style>
