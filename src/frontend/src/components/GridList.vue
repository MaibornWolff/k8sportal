<template>
  <v-layout class="ma-3">
    <v-flex xs12 md6 offset-md3>
      <v-card>
        <v-container v-bind="{ [`grid-list-${size}`]: true }" fluid>
          <v-row justify="space-between">
            <v-col md="4 text-overline mx-3">
              <h1>Services</h1>
            </v-col>
            <v-spacer></v-spacer>
            <v-col md="4" class="d-flex align-center">
              <h3>Filter:</h3>

              <v-chip-group active-class="primary--text" column>
                <v-chip
                  v-for="(chip, key) in chips"
                  :key="key"
                  x-small
                  class="ml-1"
                  @click="filterTiles(chip)"
                  >{{ chip }}</v-chip
                >
              </v-chip-group>
            </v-col>
          </v-row>
          <v-layout row wrap>
            <v-flex v-for="(tile, key) in renderTile" :key="key" xs4>
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
      chips: ["keiner"],
      selectedTiles: "all",
      filteredTiles: [],
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
    getAvailableChips() {
      this.tiles.map((tile) => {
        if (!this.chips.includes(tile.metadata.labels.chips)) {
          this.chips.push(tile.metadata.labels.chips);
        }
      });
    },
    filterTiles(chip) {
      this.isActive = !this.isActive;
      if (chip == "keiner") {
        this.selectedTiles = "all";
      } else {
        this.filteredTiles = this.tiles.filter((tile) => {
          return tile.metadata.labels.chips == chip;
        });
        this.selectedTiles = "filter";
      }
    },
  },
  mounted() {
    axios
      .get(process.env.VUE_APP_API_ENDPOINT)
      .then((res) => {
        this.tiles = res.data;

        this.getLogoNames();
        this.getImages();
        this.getAvailableChips();
      })
      .catch((err) => console.log(err));
  },
  computed: {
    renderTile() {
      return this[this.selectedTiles];
    },
    all() {
      return this.tiles;
    },
    filter() {
      return this.filteredTiles;
    },
  },
};
</script>