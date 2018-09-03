<template>
  <div class="uk-overflow-auto uk-padding-small uk-padding-remove-horizontal">
    <table class="uk-table uk-table-striped uk-table-small" v-if="rpl">
      <caption>Build Info</caption>
      <thead>
        <tr>
          <th class="uk-table-shrink">Name</th>
          <th class="uk-table-expand uk-text-left">Value</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(value, propertyName) in buildInfo" :key="propertyName">
          <td class="uk-table-shrink">{{ propertyName }}</td>
          <td class="uk-table-expand uk-text-left">{{ value }}</td>
        </tr>
      </tbody>
    </table>

    <table class="uk-table uk-table-striped uk-table-small" v-if="rpl && rpl.FileInfo.Tags.length">
      <caption>Tags</caption>
      <thead>
        <tr>
          <th class="uk-table-shrink">Name</th>
          <th class="uk-table-expand uk-text-left">Value</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="tag in rpl.FileInfo.Tags" :key="tag.Key">
          <td class="uk-table-shrink">{{ tag.Key }}</td>
          <td class="uk-table-expand uk-text-left">{{ tag.Value }}</td>
        </tr>
      </tbody>
    </table>

    <table class="uk-table uk-table-striped uk-table-small" v-if="rpl">
      <caption>Load Info</caption>
      <thead>
        <tr>
          <th class="uk-table-shrink">Name</th>
          <th class="uk-table-expand uk-text-left">Value</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(value, propertyName) in loadInfo" :key="propertyName">
          <td class="uk-table-shrink">{{ propertyName }}</td>
          <td class="uk-table-expand uk-text-left">0x{{ value.toString(16) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  computed: {
    ...mapGetters(['rpl']),
    buildInfo() {
      if (!this.rpl) {
        return [];
      }

      return {
          'Filename': this.rpl.FileInfo.Filename,
          'MinVersion': this.rpl.FileInfo.MinVersion,
          'CafeSdkVersion': this.rpl.FileInfo.CafeSdkVersion,
          'CafeSdkRevision': this.rpl.FileInfo.CafeSdkRevision,
          'Flags': '0x' + this.rpl.FileInfo.Flags.toString(16),
      }
    },
    loadInfo() {
      if (!this.rpl) {
        return [];
      }

      return {
          'TextAlign': this.rpl.FileInfo.TextAlign,
          'TextSize': this.rpl.FileInfo.TextSize,
          'DataAlign': this.rpl.FileInfo.DataAlign,
          'DataSize': this.rpl.FileInfo.DataSize,
          'LoadAlign': this.rpl.FileInfo.LoadAlign,
          'LoadSize': this.rpl.FileInfo.LoadSize,
          'TempSize': this.rpl.FileInfo.TempSize,
          'TrampAdjust': this.rpl.FileInfo.TrampAdjust,
          'SdaBase': this.rpl.FileInfo.SdaBase,
          'Sda2Base': this.rpl.FileInfo.Sda2Base,
          'TlsModuleIndex': this.rpl.FileInfo.TlsModuleIndex,
          'TlsAlignShift': this.rpl.FileInfo.TlsAlignShift,
      }
    },
  },
};
</script>

<style>
</style>
