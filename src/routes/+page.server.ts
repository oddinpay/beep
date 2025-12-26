import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async () => {
  const data = async () => {
    return (
      await fetch("https://beep-api.oddinpay.workers.dev/v1/status")
    ).json();
  };

  return data();
};
