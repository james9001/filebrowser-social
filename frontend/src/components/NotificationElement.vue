<template>
  <a :href="linkUrl">
    <div class="notification-element">
      <div
        class="notification"
        v-bind:class="{ unacknowledgedNotification: !acknowledged }"
      >
        <div class="notificationUserPicture">
          <img :src="getNotificationUserPictureImgSrc(causingUserName)" />
        </div>
        <div v-if="!acknowledged" class="new-notification-corpuscule"></div>
        <div class="notificationContent">
          <div class="notificationText">
            <span class="notificationUser">{{ causingUserName }}</span>
            {{ getNotificationText() }}
          </div>
          <div class="notificationFunctions">
            <div class="notificationTime" :title="formatIso8601AbsoluteTime(createdTime)">
              {{ formatFromNowTime(createdTime) }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </a>
</template>

<script>
import moment from "moment";
import { files as filesApi } from "@/api";

export default {
  name: "notification-element",
  props: [
    "id",
    "contextFilePath",
    "causingUserName",
    "notificationType",
    "createdTime",
    "acknowledged",
  ],
  async mounted() {
    //placeholder
  },
  beforeDestroy() {
    //placeholder
  },
  data: function () {
    //placeholder
    return {};
  },
  computed: {
    linkUrl() {
      const baseUrl = window.location.pathname.substring(
        0,
        window.location.pathname.indexOf("/files")
      );
      return encodeURI(`${baseUrl}/files${this.contextFilePath}`);
    },
  },
  methods: {
    formatFromNowTime: function (timeValue) {
      return moment(timeValue).fromNow();
    },
    formatIso8601AbsoluteTime: function (timeValue) {
      return moment(timeValue).format("YYYY-MM-DD HH:mm");
    },
    getNotificationText: function () {
      const noDirPathFileName =
        this.contextFilePath.split("/")[
          this.contextFilePath.split("/").length - 1
        ];

      let text = "";

      //TODO: i18n
      if (this.notificationType === "UserLeavesComment") {
        text += "left a comment on " + noDirPathFileName;
      }
      if (this.notificationType === "UserReactsToFile") {
        text += "reacted to " + noDirPathFileName;
      }
      if (this.notificationType === "UserReactsToComment") {
        text += "reacted to a comment on " + noDirPathFileName;
      }
      if (this.notificationType === "UserUploaded") {
        text += "uploaded something to " + this.contextFilePath;
      }

      return text;
    },
    getNotificationUserPictureImgSrc(causingUserName) {
      const imagePath = encodeURI(
        `/filebrowser-social/users/${causingUserName}.png`
      );
      return filesApi.getDownloadURL({ path: imagePath }, "thumb");
    },
  },
};
</script>
<style>
.notification {
  margin: 0em 0.5em 0.5em 0.5em;
  padding: 0.5em;
  background: var(--surfacePrimary);
  border-radius: 1em;
  border: 1px solid var(--surfacePrimary);
}

.notificationUserPicture img {
  float: left;
  width: 2em;
  height: 2em;
  background-color: black;
  border-radius: 1em;
}

.notificationContent {
  padding-left: 2.5em;
}

.notificationUser {
  color: var(--textSecondary);
}

.notificationText {
  white-space: pre-line;
}

.notificationFunctions {
  font-size: 0.8em;
}

.notificationDeleteButton {
  cursor: pointer;
  color: var(--blue);
}

.notificationTime {
  color: var(--textSecondary);
}

.unacknowledgedNotification {
  box-shadow: inset 0 0 0 1000px var(--blue_trans_05percent);
  border: 1px solid var(--blue_trans_20percent);
}

.new-notification-corpuscule {
  float: right;
  width: 0.5em;
  height: 0.5em;
  background-color: var(--blue);
  border-radius: 1em;
}
</style>
