<template>
  <div class="uk-overflow-auto uk-padding-small uk-padding-remove-horizontal">
    <table class="uk-table uk-table-striped uk-table-small" v-if="imports.functions && imports.functions.length">
      <caption>Functions</caption>
      <thead>
        <tr>
          <th class="uk-table-expand">Name</th>
          <th class="uk-table-shrink uk-text-right">RPL</th>
        </tr>
      </thead>
      <tbody>
        <RplImportRow v-for="item in imports.functions" :item="item" :key="item.name" />
      </tbody>
    </table>
    <table class="uk-table uk-table-striped uk-table-small" v-if="imports.data && imports.data.length">
      <caption>Data</caption>
      <thead>
        <tr>
          <th class="uk-table-expand">Name</th>
          <th class="uk-table-shrink uk-text-right">RPL</th>
        </tr>
      </thead>
      <tbody>
        <RplImportRow v-for="item in imports.data" :item="item" :key="item.name" />
      </tbody>
    </table>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
import RplImportRow from './RplImportRow.vue';

export default {
  components: { RplImportRow },
  computed: {
    ...mapGetters(['imports']),
  },
};
</script>

<style>
</style>
