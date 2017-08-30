<template>
  <div class="uk-background-muted uk-position-cover uk-width-1-1"
    uk-height-viewport="expand: true"
    v-show="searchTerm && searchTerm.length > 2"
  >
    <div class="uk-padding uk-margin-large-top uk-grid-divider uk-child-width-1-1 uk-grid-stack uk-overflow-auto" uk-grid>
      <div class="uk-grid-margin uk-first-column">
        <h3>Titles</h3>
        <table class="uk-table uk-table-hover">
          <tbody>
            <tr v-for="title in titleSearchResults" :key="title.ID">
              <td class="uk-table-expand uk-table-link">
                <router-link :to="{ name: 'title', params: { titleId: title.HexID }}" @click.native="clearSearch">
                  {{ title.LongNameEnglish }}
                </router-link>
              </td>
              <td class="uk-table-shrink">{{ title.PublisherEnglish }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="uk-grid-margin uk-first-column">
        <h3>Functions</h3>
        <table class="uk-table uk-table-hover">
          <tbody>
            <tr v-for="item in functionSearchResults" :key="item.ID">
              <td class="uk-table-expand uk-table-link">
                <router-link :to="{ name: 'title', params: { titleId: item.TitleHexID, rplId: item.RPLID, type: 'exports' }}" @click.native="clearSearch">
                  {{ item.Name }}
                </router-link>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="uk-grid-margin uk-first-column">
        <h3>Data</h3>
        <table class="uk-table uk-table-hover">
          <tbody>
            <tr v-for="item in dataSearchResults" :key="item.ID">
              <td class="uk-table-expand uk-table-link">
                <router-link :to="{ name: 'title', params: { titleId: item.TitleHexID, rplId: item.RPLID, type: 'exports' }}" @click.native="clearSearch">
                  {{ item.Name }}
                </router-link>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  computed: {
    ...mapGetters([
      'searchTerm',
      'titleSearchResults',
      'dataSearchResults',
      'functionSearchResults',
    ]),
  },
  methods: {
    clearSearch() {
      this.$store.dispatch('clearSearch');
    },
  },
};
</script>

<style>
</style>
