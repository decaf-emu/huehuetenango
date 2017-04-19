<template>
  <div class="uk-position-relative">
    <h3>Functions</h3>
    <div class="uk-overflow-auto">
      <table class="uk-table">
        <RplImportRow v-for="item in imports.functions" :item="item" :key="item.name" />
      </table>
    </div>

    <h3>Data</h3>
    <div class="uk-overflow-auto">
      <table class="uk-table">
        <RplImportRow v-for="item in imports.data" :item="item" :key="item.name" />
      </table>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
import RplImportRow from './RplImportRow.vue';

export default {
  components: { RplImportRow },
  props: ['titleId', 'rplId'],
  beforeMount() {
    this.listImports();
  },
  computed: {
    ...mapGetters(['imports', 'loadingImports']),
  },
  methods: {
    listImports() {
      const { titleId, rplId } = this;

      if (titleId && rplId) {
        this.$store.dispatch('listImports', { titleId, rplId });
      }
    },
  },
  watch: {
    rplId() {
      this.listImports();
    },
    titleId() {
      this.listImports();
    },
  },
};
</script>

<style>
</style>
