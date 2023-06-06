<h1 align="center">nsq_forward</h1>

<p align="center">Like nsq_to_nsq but exit after all messages are processed in source channel.</p>

<p align="center">
  <a href="https://hub.docker.com/r/iamolegga/nsq_forward">
    <img alt="Docker Image Version (latest semver)" src="https://img.shields.io/docker/v/iamolegga/nsq_forward?sort=semver">
  </a>
  <a href="https://github.com/iamolegga/nsq_forward/actions/workflows/on-push-main.yml?query=branch%3Amain">
    <img alt="GitHub Workflow Status (with branch)" src="https://img.shields.io/github/actions/workflow/status/iamolegga/nsq_forward/on-push-main.yml?branch=main">
  </a>
  <a href="https://snyk.io/test/github/iamolegga/nsq_forward">
    <img alt="Snyk Vulnerabilities for GitHub Repo (Specific Manifest)" src="https://img.shields.io/snyk/vulnerabilities/github/iamolegga/nsq_forward/go.mod" />
  </a>
  <a href="https://libraries.io/github/iamolegga/nsq_forward">
    <img alt="Libraries.io dependency status for GitHub repo" src="https://img.shields.io/librariesio/github/iamolegga/nsq_forward" />
  </a>
  <img alt="Dependabot" src="https://badgen.net/github/dependabot/iamolegga/nsq_forward" />
  <img alt="Docker Pulls" src="https://img.shields.io/docker/pulls/iamolegga/nsq_forward" />
</p>

## Usage

```shell
docker run --rm iamolegga/nsq_forward \
  -nsqd-tcp-address source-nsqd-host:4150 \
  -destination-nsqd-tcp-address destination-nsqd-host:4150 \
  -topic source_topic \
  -channel source_channel \
  -destination-topic destination_topic
```

## License

This project is licensed under the MIT License. See the LICENSE file for details.

