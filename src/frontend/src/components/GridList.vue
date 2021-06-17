<template>
  <v-layout class="ma-3">
    <v-flex xs12 sm6 offset-sm3>
      <v-card>
        <v-container v-bind="{ [`grid-list-${size}`]: true }" fluid>
          <v-row justify="space-between">
            <v-col md="4 text-overline mx-3">
              <h1>Services</h1>
            </v-col>
            <v-spacer></v-spacer>
            <!-- <v-col md="4 text-h6 mx-3">
              <span>Filter: </span>
              <span v-for="(tile, key) in tiles" :key="key">
                <v-chip v-if="tile.metadata.labels.chip">
                  {{ tile.metadata.labels.chip }}</v-chip
                >
              </span>
            </v-col> -->
          </v-row>
          <v-layout row wrap>
            <v-flex v-for="(tile, key) in tiles" :key="key" xs4>
              <v-card flat tile>
                <Tile :data="tile" />
              </v-card>
            </v-flex>
          </v-layout>
        </v-container>
      </v-card>
    </v-flex>
  </v-layout>
</template>

<script>
import Tile from "./Tile.vue";
import axios from "axios";

export default {
  name: "GridList",
  components: {
    Tile,
  },
  data() {
    return {
      size: "xl",
      tiles: [],
      images: [],
      names: [],
    };
  },
  methods: {
    // imports all image path in assets/ folder -> push in images
    importAll(r) {
      r.keys().forEach((key) =>
        this.images.push({ pathLong: r(key), pathShort: key })
      );
    },
    // get short name of all image paths -> push in names
    getLogoNames() {
      const re = /(?<=.\/)(.*?)(?=.svg)/g;
      this.importAll(require.context("../assets/", true, /\.svg$/));
      this.images.map((element) =>
        this.names.push(element.pathShort.match(re)[0])
      );
    },
    // for each tile and for each name -> if match -> insert new image attribute in tile object with image path of matched image
    // and pass via props
    getImages() {
      const img = require.context("../assets/", false, /\.svg$/);
      this.tiles.map((tile) => {
        this.names.forEach((name) => {
          if (!tile.metadata.name.includes(name)) {
            tile["image"] = img("./" + "default" + ".svg");
          }
        });

        this.names.forEach((name) => {
          if (tile.metadata.name.includes(name)) {
            tile["image"] = img("./" + name + ".svg");
          }
        });
      });
    },
  },
  mounted() {
    axios
      .get(process.env.VUE_APP_API_ENDPOINT)
      .then((res) => {
        this.tiles = res.data;

        this.getLogoNames();
        this.getImages();
      })
      .catch((err) => console.log(err));
  },
};
</script>