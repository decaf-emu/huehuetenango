<template>
  <table class="uk-table uk-table-hover" v-show="titleSearchResults && titleSearchResults.length > 0">
    <tbody>
      <tr v-for="title in titleSearchResults">
        <td class="uk-table-expand uk-table-link">
          <router-link :to="{ name: 'title', params: { titleId: title.HexID }}">
            {{ title.LongNameEnglish }}
          </router-link>
        </td>
        <td class="uk-table-shrink">{{ title.PublisherEnglish }}</td>
      </tr>
    </tbody>
  </table>
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  computed: {
    ...mapGetters(['titleSearchResults']),
  },
};
</script>

<style>
</style>
