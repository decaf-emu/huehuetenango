<template>
  <div class="uk-overflow-auto uk-padding-small uk-padding-remove-horizontal">
    <table class="uk-table uk-table-striped uk-table-small" v-if="exports.functions && exports.functions.length">
      <caption>Functions</caption>
      <thead>
        <tr>
          <th>Name</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in exports.functions" :key="item.name">
          <td>
            <code>
              {{ item }}
            </code>
          </td>
        </tr>
      </tbody>
    </table>
    <table class="uk-table uk-table-striped uk-table-small" v-if="exports.data && exports.data.length">
      <caption>Data</caption>
      <thead>
        <tr>
          <th>Name</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in exports.data" :key="item.name">
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
    ...mapGetters(['exports']),
  },
};
</script>

<style>
</style>
