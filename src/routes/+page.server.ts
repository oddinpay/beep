import type { PageServerLoad } from "./$types";
import { env } from "../env/server";

export const load: PageServerLoad = async () => {
  const data = async () => {
    return (
      await fetch(`https://${env.BEEP_HOST}/v1/status`, {
        method: "POST",
      })
    ).json();
  };

  return data();
};
