<template>
  <div class="uk-position-relative uk-padding-small">
    <h4 v-if="rpl">{{ rpl.Name }}</h4>
    <ul uk-tab>
      <li :class="{ 'uk-active': type == 'info' }">
        <router-link :to="{ name: 'title', params: { titleId, rplId, type: 'info' }}">Info</router-link>
      </li>
      <li :class="{ 'uk-active': type == 'imports' }">
        <router-link :to="{ name: 'title', params: { titleId, rplId, type: 'imports' }}">Imports</router-link>
      </li>
      <li :class="{ 'uk-active': type == 'exports' }">
        <router-link :to="{ name: 'title', params: { titleId, rplId, type: 'exports' }}">Exports</router-link>
      </li>
    </ul>

    <RplInfo v-show="type == 'info'" :titleId="titleId" :rplId="rplId" />
    <RplImports v-show="type == 'imports'" :titleId="titleId" :rplId="rplId" />
    <RplExports v-show="type == 'exports'" :titleId="titleId" :rplId="rplId" />
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import RplInfo from './RplInfo.vue';
import RplImports from './RplImports.vue';
import RplExports from './RplExports.vue';

export default {
  components: { RplInfo, RplExports, RplImports },
  props: ['titleId', 'rplId', 'type'],

  computed: {
    ...mapGetters(['rpl', 'loadingRpl']),
  },

  methods: {
    ...mapActions(['getRpl']),

    checkType() {
      const { titleId, rplId, type } = this;

      if (titleId && rplId && type !== 'imports' && type !== 'exports' && type !== 'info') {
        this.$router.replace({
          name: 'title',
          params: { titleId, rplId, type: 'info' },
        });
      }
    },

    fetchRpl() {
      const { titleId, rplId } = this;
      this.getRpl({ titleId, rplId });
    },
  },

  watch: {
    rplId() {
      this.fetchRpl();
    },

    type() {
      this.checkType();
    },
  },

  beforeMount() {
    this.fetchRpl();
    this.checkType();
  },
};
</script>

<style>
</style>
