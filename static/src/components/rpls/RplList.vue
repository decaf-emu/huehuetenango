<template>
  <ul class="uk-nav uk-nav-default">
    <li class="uk-nav-header">RPX</li>
    <li v-for="rpl in rpls" :key="rpl.ID" v-if="rpl.IsRPX">
      <router-link :to="{ name: 'title', params: { titleId, rplId: rpl.ID, type }}">{{ rpl.Name }}</router-link>
    </li>
    <li class="uk-nav-header">RPL</li>
    <li v-for="rpl in rpls" :key="rpl.ID" v-if="!rpl.IsRPX">
      <router-link :to="{ name: 'title', params: { titleId, rplId: rpl.ID, type }}">{{ rpl.Name }}</router-link>
    </li>
  </ul>
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  props: ['titleId', 'rpls', 'type', 'loading'],
};
</script>

<style>
</style>
