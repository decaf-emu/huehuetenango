<template>
  <div class="uk-position-relative">
    <h3>Functions</h3>
    <div class="uk-overflow-auto">
      <table class="uk-table">
        <tr v-for="item in exports.functions">
          <td>
            <code>
              {{ item }}
            </code>
          </td>
        </tr>
      </table>
    </div>

    <h3>Data</h3>
    <div class="uk-overflow-auto">
      <table class="uk-table">
        <tr v-for="item in exports.data">
          <td>
            <code>
              {{ item }}
            </code>
          </td>
        </tr>
      </table>
    </div>

    <div v-if="loadingExports" class="uk-overlay uk-overlay-default uk-position-cover">
      <div class="uk-position-center" uk-spinner></div>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  props: ['titleId', 'rplId'],
  beforeMount() {
    this.listExports();
  },
  computed: {
    ...mapGetters(['exports', 'loadingExports']),
  },
  methods: {
    listExports() {
      const { titleId, rplId } = this;

      if (titleId && rplId) {
        this.$store.dispatch('listExports', { titleId, rplId });
      }
    },
  },
  watch: {
    rplId() {
      this.listExports();
    },
    titleId() {
      this.listExports();
    },
  },
};
</script>

<style>
</style>
