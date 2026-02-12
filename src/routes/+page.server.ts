import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async () => {
  const data = async () => {
    return (await fetch("https://beep.oddinpay.com/v1/status")).json();
  };

  return data();
};
