<template>
  <div class="notification-pane-main noselect">
    <div class="notification-flex-container">
      <div class="notification-pane-title">
        <h2>Notifications</h2>
      </div>
      <notification-element
        v-for="notification in notifications"
        :key="notification.id"
        v-bind:id="notification.id"
        v-bind:contextFilePath="notification.contextFilePath"
        v-bind:causingUserName="notification.causingUserName"
        v-bind:notificationType="notification.notificationType"
        v-bind:createdTime="notification.createdTime"
        v-bind:acknowledged="notification.acknowledged"
      >
      </notification-element>
      <div
        class="notification-the-end"
        v-bind:class="{ notificationTheEndShown: showTheEnd }"
      >
        You've reached the end.
      </div>
      <div class="notification-load-more">
        <button
          class="button button--flat"
          v-bind:class="{ notificationTheEndHidden: showTheEnd }"
          @click="loadMore"
          title="Load More"
        >
          Load More
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { notifications as notificationsApi } from "@/api";
import NotificationElement from "@/components/NotificationElement.vue";
export default {
  name: "notification-pane",
  components: {
    NotificationElement,
  },
  props: ["filePath"],
  async mounted() {
    this.firstLoad();
  },
  beforeDestroy() {},
  data: function () {
    return {
      notifications: [],
      pageLoaded: false,
      currentlyLoading: true,
      showTheEnd: false,
      unacknowledgedNotificationCount: 0,
    };
  },
  methods: {
    async firstLoad() {
      if (!this.pageLoaded) {
        if (this.filePath) {
          await notificationsApi.acknowledgeNotificationsForFilePath(
            this.filePath
          );
        }

        this.unacknowledgedNotificationCount =
          await notificationsApi.fetchUnacknowledgedNotificationsCount();
        this.$emit(
          "update:unacknowledgedNotificationCount",
          this.unacknowledgedNotificationCount
        );

        const data = await notificationsApi.fetchNotificationPage(0);
        if (data.notifications) {
          data.notifications.forEach((notification) =>
            this.notifications.push(notification)
          );
        } else {
          this.showTheEnd = true;
        }
        this.pageLoaded = true;
        this.currentlyLoading = false;
      }
    },
    async loadMore() {
      if (!this.currentlyLoading && !this.showTheEnd) {
        this.currentlyLoading = true;
        const currentOldestNotificationId =
          this.notifications[this.notifications.length - 1].id;
        const lowId = Math.max(currentOldestNotificationId - 11, 1);
        const highId = Math.max(currentOldestNotificationId - 1, 1);
        const additionalData = await notificationsApi.fetchNotificationRange(
          lowId,
          highId
        );
        additionalData.notifications.forEach((notification) => {
          if (
            !this.notifications.find(
              (existing) => existing.id === notification.id
            )
          ) {
            this.notifications.push(notification);
          }
        });
        if (lowId == 1 && highId == 1) {
          this.showTheEnd = true;
        }
        this.currentlyLoading = false;
      }
    },
  },
};
</script>
<style>
.notification-pane-main {
  background: var(--socialContainerBg);
  border-radius: 2px;
  box-shadow: 0 0 0.25rem 0.25rem rgba(0, 0, 0, 0.2);
  overflow: auto;
  position: fixed;
  left: 0;
  right: 0;

  animation: 0.1s show forwards;
  z-index: 99999;

  min-height: calc(100%);
  max-height: calc(100%);
  overflow-y: scroll;
}

.notification-flex-container {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: left;
}

.notification-flex-container div {
  flex: 1;
}

.notification-pane-title {
  padding: 1em 1em 1em;
  text-align: center;
}

.notification-pane-title h2 {
  font-weight: 500;
  margin: 0;
}

.notification-load-more {
  text-align: center;
}

.notification-load-more input {
  width: 100%;
  text-align: center;
}

.notification-the-end {
  color: var(--textSecondary);
  text-align: center;
  padding: 1em 1em 1em;
  display: none;
}

.notificationTheEndShown {
  display: block;
}

.notificationTheEndHidden {
  display: none;
}

@media (min-width: 736px) {
  .notification-pane-main {
    left: unset;
    right: 0.25em;
    width: 25em;
    min-height: 40em;
    max-height: calc(100% - 4.25em);
    top: 4em;
  }

  .notification-pane-title {
    text-align: unset;
  }
}
</style>
