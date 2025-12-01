# Appendix C. Resources available

> …the dead hand of the academy had yet to stifle the unbridled enthusiasms of a small band of amateurs in Europe and America.
>
> — Michael D. Coe, _Breaking the Maya Code_

**Note:** The hypertext links in this appendix are to WWW sites external to the standards document. Links have been verified and updated where possible.

The resources below are mainly available from the [if-archive](https://ifarchive.org/if-archive/) maintained by David Kinder and Stephen Granade.

## Public Interpreters

Since the 1.0 Standard was first introduced, a variety of [interpreters](https://ifarchive.org/indexes/if-archiveXinfocomXinterpreters.html) have become available, for many different platforms and with a variety of features. The [ifwiki page on the Z-Machine](https://www.ifwiki.org/index.php/Z-machine#Assorted_Z-machine_interpreters) is a good place to start when looking for an interpreter.

## Testing compliance

A few story files exist to test interpreter for various Z-Machine features:

1. [Czech](https://ifarchive.org/if-archive/infocom/interpreters/tools/czech_0_8.zip) is a general Z-Machine unit test by Evin Robertson and Amir Karger.

2. [Praxix](https://eblong.com/zarf/ftp/praxix.z5) is a general Z-Machine unit test by Andrew Plotkin and Dannii Willis.

3. [TerpEtude](https://ifarchive.org/if-archive/infocom/interpreters/tools/etude.tar.Z) by Andrew Plotkin is a unit test with a focus on IO.

4. Unicode by David Kinder tests Unicode input and output.

5. [Strict Z Test](https://ifarchive.org/if-archive/infocom/interpreters/tools/strictz.z5) by Torbjorn Andersson tests an interpreter's ability to handle errors in story files.

## Compilers

Infocom's original compiler Zilch no longer exists, but Infocom's documentation for their ZIL language is available. Since the emergence of this file, there have been a few attempts at making a new ZIL compiler, most notably [ZILF](https://sourceforge.net/projects/zilf/) by Jesse McGrew.

No continuous part of the source code of any of Infocom's games is in the public domain (but see [Stu Galley's chapter of an Infocom history article](https://ifarchive.org/if-archive/infocom/articles/NZT-Zorkhistory.txt), and the IEEE article, for fragments).

Inform remains the most popular compiler for Z-Machine games. In 2006, the first public release of [Inform 7](https://inform7.com/) was made available. Inform 7 is a very different language from previous versions of Inform, with a focus on natural language and a powerful IDE.

While Inform 7 is the main focus of development, [Inform 6](https://inform-fiction.org/) is still widely used and the compiler and libraries continue to be maintained.

A few [other compilers](https://ifarchive.org/indexes/if-archiveXinfocomXcompilers.html) have been made available over the years.

## Utility programs

Mark Howell has written a toolkit of [Ztools](https://ifarchive.org/indexes/if-archiveXinfocomXtoolsXztools.html), or utility programs (1991–5, updated 1997), which includes:

1. **Txd**, a disassembler for Versions 1 to 8. (Uses the same opcode names as Inform and this document, and has an option to disassemble in Inform assembly-language syntax.)

2. **Infodump**, capable of printing the header information, object tree (with properties and attributes), dictionary and grammar tables of any Infocom or Inform-compiled game. (Understands four varieties of grammar table: Infocom pre-Version 6, Infocom Version 6, Inform GV1 and GV2.)

3. **Pix2gif**, for converting Version 6 picture data to GIF files.

4. **Check**, for verifying Infocom or Inform story files.

These continue to be maintained (by Matthew Russotto) and the first two are extremely useful. Infodump largely supersedes Mike Threepoint's vocabulary dumper Zorkword (1991–2), which was important in its day.

[Reform](https://ifarchive.org/indexes/if-archiveXinfocomXtoolsXreform.html), by Ben Rudiak-Gould, is a decompiler for Z-Machine games which will attempt to convert them into Inform 6 source code.

## Story files

1. Numerous [Inform-compiled story files](https://ifarchive.org/indexes/if-archiveXgamesXzcode.html) are publicly available.

2. A few Infocom story files are public, notably two 4-in-1 sample games (released for advertising purposes: [55.850823](https://ifarchive.org/if-archive/infocom/demos/sampler1_R55.z3) and [97.870601](https://ifarchive.org/if-archive/infocom/demos/sampler2.z3)) and [Minizork](https://ifarchive.org/if-archive/infocom/demos/minizork.z3) (a heavily abbreviated form of Zork I released with a Commodore magazine).

3. Almost all Infocom's games remain commercially available in anthologies published by [Activision](https://www.activision.com/). Copyright resides in them and they should not be available by FTP from any site.

4. A few other Infocom story files have existed but are neither on sale nor released from copyright: this applies to several of the Version 6 games, those games involving literary rights or other legal issues (Shogun, Hitch-Hiker's Guide To The Galaxy) and ephemera such as beta-test versions which have somehow passed into private circulation.

Most of the Infocom games exist in several different releases, and some were written for one Version and then ported to later ones. Zork I, for instance, has at least 11 releases: 2 early, 8 in Version 3 (with release numbers between 5 to 88 in chronological order) and one in Version 5 (release 52).

Version 1 and 2 games are extinct, though there are a few fossils in the hands of collectors.

## Documents

The definitive guide to all Infocom story files known to exist, and an indispensable reference for anyone interested in Infocom, is Paul David Doherty's [Infocom fact sheet](https://ifarchive.org/if-archive/infocom/info/fact-sheet.txt), which is regularly updated, concise and precise.

Stefan Jokisch has written a [brief specification of Infocom-format sound effects files](https://ifarchive.org/if-archive/infocom/info/sound_format.txt).

Martin Frost is the author of the [Quetzal](https://inform-fiction.org/zmachine/standards/quetzal/index.html) standard for saved-game files. Patches to adapt Zip-based interpreters to use Quetzal are available.

Andrew Plotkin is the author of the [Blorb](https://eblong.com/zarf/blorb/blorb.html) standard for packaging up images and sounds with Z-machine games.

The [Inform Technical Manual](https://inform-fiction.org/source/tm/TechMan.txt) documents the format of parsing tables used in Inform games.

The documentation for Infocom's games is available online (with permission from Activision) courtesy of the [Infocom Documentation Project](http://infodoc.plover.net/). In addition, an archive of the ["samplers"](https://ifarchive.org/indexes/if-archiveXinfocomXshipped-documentation.html) is publicly available, as is [an interesting historical archive of magazine articles](https://ifarchive.org/indexes/if-archiveXinfocomXarticles.html) concerning Infocom, and articles from Infocom's own publicity magazine.

---

_Source: [Z-Machine Standards Document - Appendix C: Resources available](https://zspec.jaredreisinger.com/C-resources)_
