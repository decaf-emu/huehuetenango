<template>
  <div class="uk-position-relative">
    <h2 v-if="rpl">{{ rpl.Name }}</h2>
    <ul uk-tab>
      <li>
        <router-link :to="{ name: 'title', params: { titleId, rplId, type: 'exports' }}">Exports</router-link>
      </li>
      <li>
        <router-link :to="{ name: 'title', params: { titleId, rplId, type: 'imports' }}">Imports</router-link>
      </li>
    </ul>

    <RplExports v-if="type == 'exports'" :titleId="titleId" :rplId="rplId" />
    <RplImports v-if="type == 'imports'" :titleId="titleId" :rplId="rplId" />
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
import RplExports from './RplExports.vue';
import RplImports from './RplImports.vue';

export default {
  components: { RplExports, RplImports },
  props: ['titleId', 'rplId', 'type'],
  beforeMount() {
    const { titleId, rplId, type } = this;

    if (titleId && rplId) {
      this.$store.dispatch('getRpl', { titleId, rplId });
    }

    this.validateType();
  },
  computed: {
    ...mapGetters(['rpl', 'loadingRpl']),
  },
  methods: {
    validateType() {
      const { titleId, rplId, type } = this;

      if (titleId && rplId && type !== 'exports' && type !== 'imports') {
        this.$router.push({
          name: 'title',
          params: { titleId, rplId, type: 'exports' },
        });
      }
    },
  },
  watch: {
    rplId(rplId) {
      const { titleId } = this;
      this.$store.dispatch('getRpl', { titleId, rplId });
    },
    type(type) {
      this.validateType();
    },
  },
};
</script>

<style>
</style>
