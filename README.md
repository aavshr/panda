## Panda

A terminal user interface for chatting with LLMs. 

I started this project because I wanted something lightweight that I can run in my terminal to talk with LLMs.
I could not keep up the project with my work and the everchanging landscape of LLMs and their capabilities. I also realized that the terminal is not the right interface for this kind of thing.

However, I had to finish this and get a v0 out. Right now it can only talk to OpenAI's API (you need an OpenAI API Key) and is very basic in features.

![Screenshot](docs/images/screenshot.png)

### Installation

#### Download Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/aavshr/panda/releases).

**macOS:**
```bash
# Extract the archive
unzip panda_<version>_Darwin_*.zip

# Remove macOS quarantine attribute (thank you Apple)
xattr -c panda

# Move to your PATH (usually /usr/local/bin but move it wherever you like in the PATH)
sudo mv panda /usr/local/bin/
```

**Linux:**
```bash
# Move to your PATH (usually /usr/local/bin but move it wherever you like in the PATH)
sudo mv panda /usr/local/bin/
```

#### Build from Source

You need to have Go installed (version 1.23 or higher).

```bash
git clone https://github.com/aavshr/panda.git
cd panda
make build
```

### Usage

Run with `panda` after installation.

**Navigation**

- `Esc` to focus out of a section
- `Enter` to focus into a section
- Use arrow keys or `hjkl` to navigate

**Chat**

- Use `Tab` to send a message

**History**

- Use `Enter` to start a new chat in that thread
- Use `Ctrl + D` to delete a thread 
- Use `/` to filter threads
