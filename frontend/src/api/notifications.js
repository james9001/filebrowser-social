import { fetchURL } from "./utils";

export async function fetchUnacknowledgedNotificationsCount() {
  const res = await fetchURL(`/api/notifications/unacknowledged`, {});

  if (res.status === 200) {
    let data = (await res.json()).Count;
    return data;
  } else {
    throw new Error(res.status);
  }
}

export async function fetchNotificationPage(pageNum) {
  const res = await fetchURL(`/api/notifications/page/${pageNum}`, {});

  if (res.status === 200) {
    let data = await res.json();
    return data;
  } else {
    throw new Error(res.status);
  }
}

export async function fetchNotificationRange(lowId, highId) {
  const res = await fetchURL(
    `/api/notifications/range?lowId=${lowId}&highId=${highId}`,
    {}
  );

  if (res.status === 200) {
    let data = await res.json();
    return data;
  } else {
    throw new Error(res.status);
  }
}

export async function acknowledgeNotificationsForFilePath(contextFilePath) {
  const res = await fetchURL(`/api/notifications/acknowledge`, {
    method: "POST",
    body: JSON.stringify({
      data: {
        contextFilePath: contextFilePath,
      },
      which: [],
      what: "acknowledge_notification",
    }),
  });

  if (res.status !== 200) {
    throw new Error(res.status);
  }
}

export async function userUploadedNotificationPost(contextFilePath) {
  const res = await fetchURL(`/api/notifications/useruploaded`, {
    method: "POST",
    body: JSON.stringify({
      data: {
        contextFilePath: contextFilePath,
      },
      which: [],
      what: "user_uploaded_notification",
    }),
  });

  if (res.status !== 200) {
    throw new Error(res.status);
  }
}
