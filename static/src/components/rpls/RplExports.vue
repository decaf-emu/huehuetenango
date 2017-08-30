<template>
  <div class="uk-position-relative">
    <h5>Functions</h5>
    <div class="uk-overflow-auto">
      <table class="uk-table">
        <tr v-for="item in exports.functions" :key="item.name">
          <td>
            <code>
              {{ item }}
            </code>
          </td>
        </tr>
      </table>
    </div>
    <h5>Data</h5>
    <div class="uk-overflow-auto">
      <table class="uk-table">
        <tr v-for="item in exports.data" :key="item.name">
          <td>
            <code>
              {{ item }}
            </code>
          </td>
        </tr>
      </table>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  computed: {
    ...mapGetters(['exports']),
  },
};
</script>

<style>
</style>
