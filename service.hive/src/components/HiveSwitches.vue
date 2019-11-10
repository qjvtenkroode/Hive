<template>
  <div>
    <button @click="getAssets" class="text-orange-500 hover:text-white hover:bg-orange-500 border border-orange-500 text-xs font-semibold rounded-full px-4 py-1 leading-normal">Refresh assets</button>
    <div class="flex flex-wrap justify-center">
      <HiveSwitch v-for="asset in assets" v-bind:key="asset.identifier" :asset="asset"></HiveSwitch>
    </div>
  </div>
</template>

<script>
  import HiveSwitch from './HiveSwitch.vue'

  export default {
    name: 'HiveSwitches',

    components: {
      HiveSwitch
    },

    mounted: function() {
      this.getAssets();
    },

    data: function() {
      return {
        assets: [],
      }
    },

    methods: {
      getAssets: function() {
        axios.get('http://' + host + ':8000/assets/')
          .then(response => this.assets = response.data);
      }
    }
    /*,

    computed: {
      shelly: function() {
        return this.assets.filter(function(u) {
          return u.type == "shelly1" || u.type == "shelly1pm"
        })
      }
    }*/
  }
</script>

<style>

</style>
