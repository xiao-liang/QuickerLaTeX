# LaTeX to (Wordpress) HTML Converter

This is a command-line tool in [Golang](https://golang.org/) that converts LaTeX (in its restricted form) to HTML codes for Wordpress posts. Note that this tool is supposed to be used in combination with the Wordpress plugin [QuickLaTeX](http://www.holoborodko.com/pavel/quicklatex/) (Yes, an extension of an extention). Though it does not work by itself, it will ease your pain in converting LaTeX codes into html even with the QuickLaTex plugin. That is why it is named **QuickerLaTeX**. 

## Features (and the motivation)
Let us first take a look at what it does and what it does not do. All the supported LaTeX features can be put into to the following two categories:
#### The ones supported by QuickLaTeX:
[QuickLaTeX](http://www.holoborodko.com/pavel/quicklatex/)  already supports a decent amount of LaTeX features, most notably, including:
  * `equation` environment with optional labeling `\label{_label_}` and cross-referencing `\ref{_label_}`
  * `align` environment. (this is AMAZING!)  

For a complete list, visit the [webpage of Pavel Holoborodko](http://www.holoborodko.com/pavel/quicklatex/) (the amazing developer)

These features are actually already enough for many Wordpress bloggers who use an moderate amount of LaTeX codes.

For these features, our **QuickerLaTeX** tool does nothing but simply keep the LaTeX codes as it is, to allow QuickLaTeX to deal with it.



#### Features in Addition to QuickLaTeX:
  * hyperref urls via `\href{_url_}{_work_}`; 
  * sectioning commands `\section{_title_}` and `\subsection{_title_}`. Optional `\label{section:_label_}` and `\ref{section:_label_}` are supported;
  * environments for theorems and lemmata: `\begin{theorem}...\end{theorem}`, `\begin{lemma}...\end{lemma}`. Similar as the above, optional labeling and cross-referencing are supported; 
  * the `itemize` and `enumerate` environments in their simplest forms (taking no parameters)
  * macros in its simplest format, i.e. `\newcommand{\_symbol_}{_value_}` without optional parameters


> **Some words on the motivation:** the features listed in this category are the reason why I build this tool. There are several other similar tools (e.g., the [LaTeX2WP](https://lucatrevisan.wordpress.com/latex-to-wordpress/) by Prof. Luca Trevisan). But they does not suit my need very well. For example,  LaTeX2WP's output does not always align well in-line, and it seems to depend on Jetpack's "[Beautiful Math with LaTeX](https://jetpack.com/support/beautiful-math-with-latex/)", which may result in low-quality pictures with certain blog themes (thus, not "beautiful" any more). The latter may be fixed by doing some extra configurations (not sure?), which I found harder than writing a converter myself. Then, **QuickerLaTeX**, here it is.


## Usage
1. Install and active  [QuickLaTeX](http://www.holoborodko.com/pavel/quicklatex/). (You may want to see instructions for installing and usage by clicking the link.)
2. Write your code in your favorite latex editor. There are certain syntax that should be followed. See [The Rules](#rules) below. Assume the file is `example.tex`
3. Install [Golang](https://golang.org/) on your machine and make sure that it works (e.g., GOPATH and GOROOT are configured properly). There are thousands of tutorials on-line on this topic for all kinds of operating systems. Thus, I assume that Golang is already installed on you machine properly.
4. Copy the `main.go` file in `QuickerLaTeX\src\` folder to the same directory of the `example.tex` file
5. Open your terminal and `cd` into the the same directory of the `example.tex` file. (Note that and the `mian.go` is also in this folder by now.)
6. Execute `go run main.go example.tex`. Then, a file called `example.txt` will be generated.
7. Copy the content of `example.txt` to the Wordpress post editor in "**code view**" mode. Namely, you need to past the content there as plain html code, instead of using the **visual editor**. 


## The Rules <a name="rules"></a>
* Only use LaTeX features that are explicitly listed as "supported" by [QuickLaTeX](http://www.holoborodko.com/pavel/quicklatex/) or **QuickerLaTeX** (the current tool).
* Obey the rules for [QuickLaTeX](http://www.holoborodko.com/pavel/quicklatex/). (Indeed, my tool is just an extension of QuickLaTeX).
* Labels for both sections and subsections should start with `section:`. Namely, it should be of the form`\label{section:_your_label_}`
* Similarly, labeling `theorem` environment with  `\label{theorem:_your_label_}`, and `lemma` environment with  `\label{lemma:_your_label_}`
* As mentioned earlier, `\itemize`, `\enumerate` and `\newcommand` are supported in their simplest form.

## Acknowledgements
This work is inspired by the [LaTeX2WP](https://lucatrevisan.wordpress.com/latex-to-wordpress/) of Prof. Luca Trevisan, together with the [modified version](https://github.com/seaneberhard/latex2wp) by Sean Eberhard. 

## The Final Note
The code is not clean as I write it in a huge hurry (roughly 1 day and with half of it reading Prof. Trevisan's code), with the hope to make it work ASAP. The code is not optimized. For several features, I'm sure there are better ways to achieved that than the one implemented by me. Also, exceptions and errors are not handled. I might come back to optimize the code, and support additional features when I really need it. 