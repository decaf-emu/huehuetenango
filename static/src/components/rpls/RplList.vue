<template>
  <table class="uk-table uk-table-hover">
    <tr v-for="rpl in rpls">
      <td class="uk-table-link">
        <router-link :to="{ name: 'title', params: { titleId, rplId: rpl.ID, type }}">{{ rpl.Name }}</router-link>
      </td>
    </tr>
  </table>
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  props: ['titleId', 'rpls', 'type', 'loading'],
};
</script>

<style>
</style>
