import http from 'k6/http';
import { sleep } from 'k6';
import * as faker from 'faker/locale/en_US';
import { Rate } from 'k6/metrics';

const baseUrl = 'http://localhost:8080'

const deliveredFailedRate = new Rate('failed delivered calls');
const bouncedFailedRate = new Rate('failed bounced calls');
const getDomainFailedRate = new Rate('failed get domain calls');

export const options = {
  vus: 100,
  duration: '60s',
  thresholds: {
    'failed delivered calls': ['rate<0.1'],
    'failed bounced calls': ['rate<0.1'],
    'failed get domain calls': ['rate<0.1'],
    'http_req_duration': ['p(95)<200']
  }
};

/**
 * Sets up the actual functions to call to load tests.
 */
export default function(data) {
  delivered();
  bounced();
  getDomain();
  sleep(1);
}

/**
 * Calls the delivered endpoint.
 */
const delivered = () => {
  const result = http.put(`${baseUrl}/events/${randomDomain()}/delivered`, '', {
    tags: { name: 'Delivered'}
  });
  deliveredFailedRate.add(result.status != 200);
}

/**
 * Calls the bounced endpoint.
 */
const bounced = () => {
  const result = http.put(`${baseUrl}/events/${randomDomain()}/bounced`, '', {
    tags: { name: 'Bounced'}
  });

  bouncedFailedRate.add(result.status != 200);
}

/**
 * Calls the get domain endpoint.
 */
const getDomain = () => {
  let url = `${baseUrl}/domains/${randomDomain()}`;
  console.log(url);
  const result = http.get(url, {
    tags: { name: 'GetDomain'}
  });

  getDomainFailedRate.add(result.status != 200);
}


/**
 * Returns a random domain with a probability of it being one already visited.
 */
const randomDomain = () => {
  return faker.internet.domainName();
}
