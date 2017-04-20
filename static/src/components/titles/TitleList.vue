<template>
  <div class="uk-container uk-container-expand uk-position-relative">
    <table class="uk-table uk-table-hover">
      <thead>
        <tr>
          <th class="uk-table-expand">Name</th>
          <th class="uk-table-shrink">Short Name</th>
          <th class="uk-table-shrink">Publisher</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="title in allTitles">
          <td class="uk-table-link">
            <router-link :to="{ name: 'title', params: { titleId: title.HexID }}">
              {{ title.LongNameEnglish }}
            </router-link>
          </td>
          <td class="uk-table-link">
            <router-link :to="{ name: 'title', params: { titleId: title.HexID }}">
              {{ title.ShortNameEnglish }}
            </router-link>
          </td>
          <td>{{ title.PublisherEnglish }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';

export default {
  computed: {
    ...mapGetters(['allTitles', 'loadingAllTitles']),
  },

  methods: {
    ...mapActions([
      'getAllTitles',
      'clearTitle',
      'clearTitleRpls',
      'clearRpl',
      'clearImports',
      'clearExports',
    ]),
  },

  head: {
    title: {
      inner: 'Titles',
    },
  },

  beforeMount() {
    this.getAllTitles();

    this.clearTitle();
    this.clearTitleRpls();
    this.clearRpl();
    this.clearImports();
    this.clearExports();
  },
};
</script>

<style>
</style>
