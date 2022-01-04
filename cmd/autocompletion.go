package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"os"
)

var compCmd = cobra.Command{
	Use:   "completion bash",
	Short: "generate autocompletion script for bash",
	Long: `To load completion run

. <(qshell completion <bash|zsh>)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile or ~/.zshrc
. <(qshell completion <bash|zsh>)
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, params []string) {
		shName := params[0]

		switch shName {
		case "bash":
			rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			zsh_head := "#compdef qshell\n"

			os.Stdout.Write([]byte(zsh_head))

			zsh_initialization := `
__qshell_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand
	source "$@"
}
__qshell_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift
		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__qshell_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}
__qshell_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?
	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}
__qshell_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}
__qshell_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}
__qshell_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}
__qshell_filedir() {
	local RET OLD_IFS w qw
	__qshell_debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi
	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"
	IFS="," __qshell_debug "RET=${RET[@]} len=${#RET[@]}"
	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__qshell_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}
__qshell_quote() {
    if [[ $1 == \'* || $1 == \"* ]]; then
        # Leave out first character
        printf %q "${1:1}"
    else
	printf %q "$1"
    fi
}
autoload -U +X bashcompinit && bashcompinit
# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi
__qshell_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__qshell_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__qshell_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__qshell_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__qshell_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__qshell_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/builtin declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__qshell_type/g" \
	<<'BASH_COMPLETION_EOF'
`
			os.Stdout.Write([]byte(zsh_initialization))

			buf := new(bytes.Buffer)
			rootCmd.GenBashCompletion(buf)
			os.Stdout.Write(buf.Bytes())

			zsh_tail := `
BASH_COMPLETION_EOF
}
__qshell_bash_source <(__qshell_convert_bash_to_zsh)
_complete qshell 2>/dev/null
`
			os.Stdout.Write([]byte(zsh_tail))
		}
	},
}

func init() {
	rootCmd.AddCommand(&compCmd)
}
