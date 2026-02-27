import type { PageServerLoad } from "./$types";
import { env } from "../env/server";

export const load: PageServerLoad = async () => {
  const data = async () => {
    return (await fetch(`https://${env.ODDIN_HOST}/v1/status`)).json();
  };

  return data();
};
